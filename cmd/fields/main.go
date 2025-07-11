package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/hashicorp/go-version"
	"github.com/iancoleman/strcase"
)

type replacement struct {
	Old string
	New string
}

var fieldReps = []replacement{
	{"Dhcpdv6", "DHCPDV6"},

	{"Dhcpd", "DHCPD"},
	{"Idx", "IDX"},
	{"Ipsec", "IPSec"},
	{"Ipv6", "IPV6"},
	{"Openvpn", "OpenVPN"},
	{"Tftp", "TFTP"},
	{"Wlangroup", "WLANGroup"},

	{"Bc", "Broadcast"},
	{"Dhcp", "DHCP"},
	{"Dns", "DNS"},
	{"Dpi", "DPI"},
	{"Dtim", "DTIM"},
	{"Firewallgroup", "FirewallGroup"},
	{"Fixedip", "FixedIP"},
	{"Icmp", "ICMP"},
	{"Id", "ID"},
	{"Igmp", "IGMP"},
	{"Ip", "IP"},
	{"Leasetime", "LeaseTime"},
	{"Mac", "MAC"},
	{"Mcastenhance", "MulticastEnhance"},
	{"Minrssi", "MinRSSI"},
	{"Monthdays", "MonthDays"},
	{"Nat", "NAT"},
	{"Networkconf", "Network"},
	{"Networkgroup", "NetworkGroup"},
	{"Pd", "PD"},
	{"Pmf", "PMF"},
	{"Portconf", "PortProfile"},
	{"Qos", "QOS"},
	{"Radiusprofile", "RADIUSProfile"},
	{"Radius", "RADIUS"},
	{"Ssid", "SSID"},
	{"Startdate", "StartDate"},
	{"Starttime", "StartTime"},
	{"Stopdate", "StopDate"},
	{"Stoptime", "StopTime"},
	{"Tcp", "TCP"},
	{"Udp", "UDP"},
	{"Usergroup", "UserGroup"},
	{"Utc", "UTC"},
	{"Vlan", "VLAN"},
	{"Vpn", "VPN"},
	{"Wan", "WAN"},
	{"Wep", "WEP"},
	{"Wlan", "WLAN"},
	{"Wpa", "WPA"},
}

var fileReps = []replacement{
	{"WlanConf", "WLAN"},
	{"Dhcp", "DHCP"},
	{"Wlan", "WLAN"},
	{"NetworkConf", "Network"},
	{"PortConf", "PortProfile"},
	{"RadiusProfile", "RADIUSProfile"},
	{"ApGroups", "APGroup"},
	{"DnsRecord", "DNSRecord"},
}

type Resource struct {
	StructName     string
	ResourcePath   string
	Types          map[string]*FieldInfo
	FieldProcessor func(name string, f *FieldInfo) error
}

type FieldInfo struct {
	FieldName           string
	JSONName            string
	FieldType           string
	FieldValidation     string
	OmitEmpty           bool
	IsArray             bool
	Fields              map[string]*FieldInfo
	CustomUnmarshalType string
	CustomUnmarshalFunc string
}

func NewResource(structName string, resourcePath string) *Resource {
	baseType := NewFieldInfo(structName, resourcePath, "struct", "", false, false, "")
	resource := &Resource{
		StructName:   structName,
		ResourcePath: resourcePath,
		Types: map[string]*FieldInfo{
			structName: baseType,
		},
		FieldProcessor: func(name string, f *FieldInfo) error { return nil },
	}

	// Since template files iterate through map keys in sorted order, these initial fields
	// are named such that they stay at the top for consistency. The spacer items create a
	// blank line in the resulting generated file.
	//
	// This hack is here for stability of the generatd code, but can be removed if desired.
	baseType.Fields = map[string]*FieldInfo{
		"   ID":      NewFieldInfo("ID", "_id", "string", "", true, false, ""),
		"   SiteID":  NewFieldInfo("SiteID", "site_id", "string", "", true, false, ""),
		"   _Spacer": nil,

		"  Hidden":   NewFieldInfo("Hidden", "attr_hidden", "bool", "", true, false, ""),
		"  HiddenID": NewFieldInfo("HiddenID", "attr_hidden_id", "string", "", true, false, ""),
		"  NoDelete": NewFieldInfo("NoDelete", "attr_no_delete", "bool", "", true, false, ""),
		"  NoEdit":   NewFieldInfo("NoEdit", "attr_no_edit", "bool", "", true, false, ""),
		"  _Spacer":  nil,

		" _Spacer": nil,
	}

	switch {
	case resource.IsSetting():
		resource.ResourcePath = strcase.ToSnake(strings.TrimPrefix(structName, "Setting"))
		baseType.Fields[" Key"] = NewFieldInfo("Key", "key", "string", "", false, false, "")

		if resource.StructName == "SettingUsg" {
			// Removed in v7, retaining for backwards compatibility
			baseType.Fields["MdnsEnabled"] = NewFieldInfo(
				"MdnsEnabled",
				"mdns_enabled",
				"bool",
				"",
				false,
				false,
				"",
			)
		}
	case resource.StructName == "DNSRecord":
		resource.ResourcePath = "static-dns"
	case resource.StructName == "Device":
		baseType.Fields[" MAC"] = NewFieldInfo("MAC", "mac", "string", "", true, false, "")
		baseType.Fields["Adopted"] = NewFieldInfo(
			"Adopted",
			"adopted",
			"bool",
			"",
			false,
			false,
			"",
		)
		baseType.Fields["Model"] = NewFieldInfo("Model", "model", "string", "", true, false, "")
		baseType.Fields["State"] = NewFieldInfo(
			"State",
			"state",
			"DeviceState",
			"",
			false,
			false,
			"",
		)
		baseType.Fields["Type"] = NewFieldInfo("Type", "type", "string", "", true, false, "")
	case resource.StructName == "User":
		baseType.Fields[" IP"] = NewFieldInfo(
			"IP",
			"ip",
			"string",
			"non-generated field",
			true,
			false,
			"",
		)
		baseType.Fields[" DevIdOverride"] = NewFieldInfo(
			"DevIdOverride",
			"dev_id_override",
			"int",
			"non-generated field",
			true,
			false,
			"",
		)
	case resource.StructName == "WLAN":
		// this field removed in v6, retaining for backwards compatibility
		baseType.Fields["WLANGroupID"] = NewFieldInfo(
			"WLANGroupID",
			"wlangroup_id",
			"string",
			"",
			false,
			false,
			"",
		)
	}

	return resource
}

func NewFieldInfo(
	fieldName string,
	jsonName string,
	fieldType string,
	fieldValidation string,
	omitempty bool,
	isArray bool,
	customUnmarshalType string,
) *FieldInfo {
	return &FieldInfo{
		FieldName:           fieldName,
		JSONName:            jsonName,
		FieldType:           fieldType,
		FieldValidation:     fieldValidation,
		OmitEmpty:           omitempty,
		IsArray:             isArray,
		CustomUnmarshalType: customUnmarshalType,
	}
}

func cleanName(name string, reps []replacement) string {
	for _, rep := range reps {
		name = strings.ReplaceAll(name, rep.Old, rep.New)
	}

	return name
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] version\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	versionBaseDirFlag := flag.String(
		"version-base-dir",
		"fields",
		"The base directory for version JSON files",
	)
	outputDirFlag := flag.String(
		"output-dir",
		"unifi",
		"The output directory of the generated Go code",
	)
	downloadOnly := flag.Bool(
		"download-only",
		false,
		"Only download and build the fields JSON directory, do not generate",
	)
	useLatestVersion := flag.Bool("latest", true, "Use the latest available version")

	flag.Parse()

	specifiedVersion := flag.Arg(0)
	if specifiedVersion != "" && *useLatestVersion {
		fmt.Print("error: cannot specify version with latest\n\n")
		usage()
		os.Exit(1)
	} else if specifiedVersion == "" && !*useLatestVersion {
		fmt.Print("error: must specify version or latest\n\n")
		usage()
		os.Exit(1)
	}

	var unifiVersion *version.Version
	var unifiDownloadUrl *url.URL
	var err error

	if *useLatestVersion {
		unifiVersion, unifiDownloadUrl, err = latestUnifiVersion()
		if err != nil {
			panic(err)
		}
	} else {
		unifiVersion, err = version.NewVersion(specifiedVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		unifiDownloadUrl, err = url.Parse(fmt.Sprintf("https://dl.ui.com/unifi/%s/unifi_sysvinit_all.deb", unifiVersion))
		if err != nil {
			panic(err)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fieldsDir := filepath.Join(wd, *versionBaseDirFlag, fmt.Sprintf("v%s", unifiVersion))
	outDir := filepath.Join(wd, *outputDirFlag)

	fieldsInfo, err := os.Stat(fieldsDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		err = os.MkdirAll(fieldsDir, 0o755)
		if err != nil {
			panic(err)
		}

		// download fields, create
		jarFile, err := downloadJar(unifiDownloadUrl, fieldsDir)
		if err != nil {
			panic(err)
		}

		err = extractJSON(jarFile, fieldsDir)
		if err != nil {
			panic(err)
		}

		err = copyCustom(fieldsDir)
		if err != nil {
			panic(err)
		}

		fieldsInfo, err = os.Stat(fieldsDir)
		if err != nil {
			panic(err)
		}
	}
	if !fieldsInfo.IsDir() {
		panic("version info isn't a directory")
	}

	if *downloadOnly {
		fmt.Println("Fields JSON ready!")
		os.Exit(0)
	}

	fieldsFiles, err := os.ReadDir(fieldsDir)
	if err != nil {
		panic(err)
	}

	for _, fieldsFile := range fieldsFiles {
		name := fieldsFile.Name()
		ext := filepath.Ext(name)

		switch name {
		case "AuthenticationRequest.json", "Setting.json", "Wall.json":
			continue
		}

		if filepath.Ext(name) != ".json" {
			continue
		}

		name = name[:len(name)-len(ext)]

		urlPath := strings.ToLower(name)
		structName := cleanName(name, fileReps)

		goFile := strcase.ToSnake(structName) + ".generated.go"
		fieldsFilePath := filepath.Join(fieldsDir, fieldsFile.Name())
		b, err := os.ReadFile(fieldsFilePath)
		if err != nil {
			fmt.Printf("skipping file %s: %s", fieldsFile.Name(), err)
			continue
		}

		resource := NewResource(structName, urlPath)

		switch resource.StructName {
		case "Account":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				switch name {
				case "IP", "NetworkID":
					f.OmitEmpty = true
				}
				return nil
			}
		case "ChannelPlan":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				switch name {
				case "Channel", "BackupChannel", "TxPower":
					if f.FieldType == "string" {
						f.CustomUnmarshalType = "numberOrString"
					}
				}
				return nil
			}
		case "Device":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				switch name {
				case "X", "Y":
					f.FieldType = "float64"
				case "StpPriority":
					f.FieldType = "string"
					f.CustomUnmarshalType = "numberOrString"
				case "Ht":
					f.FieldType = "int"
				case "Channel", "BackupChannel", "TxPower":
					if f.FieldType == "string" {
						f.CustomUnmarshalType = "numberOrString"
					}
				case "LteExtAnt", "LtePoe":
					f.CustomUnmarshalType = "booleanishString"
				}

				f.OmitEmpty = true
				switch name {
				case "PortOverrides":
					f.OmitEmpty = false
				}

				return nil
			}
		case "Network":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				switch name {
				case "InternetAccessEnabled", "IntraNetworkAccessEnabled":
					if f.FieldType == "bool" {
						f.CustomUnmarshalType = "*bool"
						f.CustomUnmarshalFunc = "emptyBoolToTrue"
					}
				case "IPSecEspLifetime", "IPSecIkeLifetime":
					f.FieldType = "*int"
				}
				return nil
			}
		case "SettingGlobalAp":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				if strings.HasPrefix(name, "6E") {
					f.FieldName = strings.Replace(f.FieldName, "6E", "SixE", 1)
				}

				return nil
			}
		case "SettingMgmt":
			sshKeyField := NewFieldInfo(
				resource.StructName+"XSshKeys",
				"x_ssh_keys",
				"struct",
				"",
				false,
				false,
				"",
			)
			sshKeyField.Fields = map[string]*FieldInfo{
				"name":    NewFieldInfo("Name", "name", "string", "", false, false, ""),
				"keyType": NewFieldInfo("KeyType", "type", "string", "", false, false, ""),
				"key":     NewFieldInfo("Key", "key", "string", "", false, false, ""),
				"comment": NewFieldInfo("Comment", "comment", "string", "", false, false, ""),
				"date":    NewFieldInfo("Date", "date", "string", "", false, false, ""),
				"fingerprint": NewFieldInfo(
					"Fingerprint",
					"fingerprint",
					"string",
					"",
					false,
					false,
					"",
				),
			}
			resource.Types[sshKeyField.FieldName] = sshKeyField

			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				if name == "XSshKeys" {
					f.FieldType = sshKeyField.FieldName
				}
				return nil
			}
		case "SettingUsg":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				if strings.HasSuffix(name, "Timeout") && name != "ArpCacheTimeout" {
					f.FieldType = "int"
					f.CustomUnmarshalType = "emptyStringInt"
				}
				return nil
			}
		case "User":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				switch name {
				case "Blocked":
					f.FieldType = "bool"
				case "LastSeen":
					f.FieldType = "int"
					f.CustomUnmarshalType = "emptyStringInt"
				}
				return nil
			}
		case "WLAN":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				switch name {
				case "ScheduleWithDuration":
					// always send schedule, so we can empty it if we want to
					f.OmitEmpty = false
				}
				return nil
			}
		case "DNSRecord":
			resource.FieldProcessor = func(name string, f *FieldInfo) error {
				switch name {
				case "Hidden", "NoDelete", "NoEdit", "Enabled":
					f.FieldType = "bool"
				case "Priority", "Ttl", "Weight":
					f.FieldType = "int"
					f.CustomUnmarshalType = "emptyStringInt"
				}
				return nil
			}
		}

		err = resource.processJSON(b)
		if err != nil {
			fmt.Printf("skipping file %s: %s", fieldsFile.Name(), err)
			continue
		}

		var code string
		if code, err = resource.generateCode(false); err != nil {
			panic(err)
		}

		_ = os.Remove(filepath.Join(outDir, goFile))
		if err := os.WriteFile(filepath.Join(outDir, goFile), ([]byte)(code), 0o644); err != nil {
			panic(err)
		}

		implFile := strcase.ToSnake(structName) + ".go"
		implFilePath := filepath.Join(outDir, implFile)

		if _, err := os.Stat(implFilePath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				var implCode string
				if implCode, err = resource.generateCode(true); err != nil {
					panic(err)
				}

				if err := os.WriteFile(filepath.Join(implFilePath), ([]byte)(implCode), 0o644); err != nil {
					panic(err)
				}
			}
		}
	}

	// Write version file.
	versionGo := fmt.Appendf(nil, `
// Generated code. DO NOT EDIT.

package unifi

const UnifiVersion = %q
`, unifiVersion)

	versionGo, err = format.Source(versionGo)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(outDir, "version.generated.go"), versionGo, 0o644); err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", outDir)
}

func (r *Resource) IsSetting() bool {
	return strings.HasPrefix(r.StructName, "Setting")
}

func (r *Resource) IsDevice() bool {
	return r.StructName == "Device"
}

func (r *Resource) IsV2() bool {
	return slices.Contains([]string{
		"DNSRecord",
		"ApGroup",
	}, r.StructName)
}

func (r *Resource) processFields(fields map[string]any) {
	t := r.Types[r.StructName]
	for name, validation := range fields {
		fieldInfo, err := r.fieldInfoFromValidation(name, validation)
		if err != nil {
			continue
		}

		t.Fields[fieldInfo.FieldName] = fieldInfo
	}
}

func (r *Resource) fieldInfoFromValidation(name string, validation any) (*FieldInfo, error) {
	fieldName := strcase.ToCamel(name)
	fieldName = cleanName(fieldName, fieldReps)

	empty := &FieldInfo{}
	var fieldInfo *FieldInfo

	switch validation := validation.(type) {
	case []any:
		if len(validation) == 0 {
			fieldInfo = NewFieldInfo(fieldName, name, "string", "", false, true, "")
			err := r.FieldProcessor(fieldName, fieldInfo)
			return fieldInfo, err
		}
		if len(validation) > 1 {
			return empty, fmt.Errorf("unknown validation %#v", validation)
		}

		fieldInfo, err := r.fieldInfoFromValidation(name, validation[0])
		if err != nil {
			return empty, err
		}

		fieldInfo.OmitEmpty = true
		fieldInfo.IsArray = true

		err = r.FieldProcessor(fieldName, fieldInfo)
		return fieldInfo, err

	case map[string]any:
		typeName := r.StructName + fieldName

		result := NewFieldInfo(fieldName, name, typeName, "", true, false, "")
		result.Fields = make(map[string]*FieldInfo)

		for name, fv := range validation {
			child, err := r.fieldInfoFromValidation(name, fv)
			if err != nil {
				return empty, err
			}

			result.Fields[child.FieldName] = child
		}

		err := r.FieldProcessor(fieldName, result)
		r.Types[typeName] = result
		return result, err

	case string:
		fieldValidation := validation
		normalized := normalizeValidation(validation)

		omitEmpty := false

		switch normalized {
		case "falsetrue", "truefalse":
			fieldInfo = NewFieldInfo(fieldName, name, "bool", "", omitEmpty, false, "")
			return fieldInfo, r.FieldProcessor(fieldName, fieldInfo)
		default:
			if _, err := strconv.ParseFloat(normalized, 64); err == nil {
				if normalized == "09" || normalized == "09.09" {
					fieldValidation = ""
				}

				if strings.Contains(normalized, ".") {
					if strings.Contains(validation, "\\.){3}") {
						break
					}

					omitEmpty = true
					fieldInfo = NewFieldInfo(fieldName, name, "float64", fieldValidation, omitEmpty, false, "")
					return fieldInfo, r.FieldProcessor(fieldName, fieldInfo)
				}

				omitEmpty = true
				fieldInfo = NewFieldInfo(fieldName, name, "int", fieldValidation, omitEmpty, false, "")
				fieldInfo.CustomUnmarshalType = "emptyStringInt"
				return fieldInfo, r.FieldProcessor(fieldName, fieldInfo)
			}
		}
		if validation != "" && normalized != "" {
			fmt.Printf("normalize %q to %q\n", validation, normalized)
		}

		omitEmpty = omitEmpty || (!strings.Contains(validation, "^$") && !strings.HasSuffix(fieldName, "Id"))
		fieldInfo = NewFieldInfo(fieldName, name, "string", fieldValidation, omitEmpty, false, "")
		return fieldInfo, r.FieldProcessor(fieldName, fieldInfo)
	}

	return empty, fmt.Errorf("unable to determine type from validation %q", validation)
}

func (r *Resource) processJSON(b []byte) error {
	var fields map[string]any
	err := json.Unmarshal(b, &fields)
	if err != nil {
		return err
	}

	r.processFields(fields)

	return nil
}

//go:embed api.go.tmpl
var apiGoTemplate string

//go:embed client.go.tmpl
var clientGoTemplate string

func (r *Resource) generateCode(isImpl bool) (string, error) {
	var err error
	var buf bytes.Buffer
	writer := io.Writer(&buf)

	var tpl *template.Template
	if isImpl {
		tpl = template.Must(template.New("client.go.tmpl").Parse(clientGoTemplate))
	} else {
		tpl = template.Must(template.New("api.go.tmpl").Parse(apiGoTemplate))
	}

	err = tpl.Execute(writer, r)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format source: %w", err)
	}

	return string(src), err
}

func normalizeValidation(re string) string {
	re = strings.ReplaceAll(re, "\\d", "[0-9]")
	re = strings.ReplaceAll(re, "[-+]?", "")
	re = strings.ReplaceAll(re, "[+-]?", "")
	re = strings.ReplaceAll(re, "[-]?", "")
	re = strings.ReplaceAll(re, "\\.", ".")
	re = strings.ReplaceAll(re, "[.]?", ".")

	quants := regexp.MustCompile(`\{\d*,?\d*\}|\*|\+|\?`)
	re = quants.ReplaceAllString(re, "")

	control := regexp.MustCompile(`[\(\[\]\)\|\-\$\^]`)
	re = control.ReplaceAllString(re, "")

	re = strings.TrimPrefix(re, "^")
	re = strings.TrimSuffix(re, "$")

	return re
}
