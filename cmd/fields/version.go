package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-version"
)

func latestUnifiVersion() (*version.Version, *url.URL, error) {
	url, err := url.Parse(firmwareUpdateApi)
	if err != nil {
		return nil, nil, err
	}

	query := url.Query()
	query.Add("filter", firmwareUpdateApiFilter("eq", "channel", releaseChannel))
	query.Add("filter", firmwareUpdateApiFilter("eq", "product", unifiControllerProduct))
	query.Add("filter", firmwareUpdateApiFilter("lt", "version", maxVersion))
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var respData firmwareUpdateApiResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, nil, err
	}

	for _, firmware := range respData.Embedded.Firmware {
		if firmware.Platform != debianPlatform {
			continue
		}

		return firmware.Version.Core(), firmware.Links.Data.Href, nil
	}

	return nil, nil, nil
}

// latestUnifiOSVersion queries the firmware update API for the latest
// unifi-os-server release on linux-x64. Unlike latestUnifiVersion, it
// has no maxVersion cap (10.x is always reached via this path). Returns
// the OS version (e.g. v5.1.21) and the download URL.
func latestUnifiOSVersion() (*version.Version, *url.URL, error) {
	u, err := url.Parse(firmwareUpdateApi)
	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Add("filter", firmwareUpdateApiFilter("eq", "channel", releaseChannel))
	q.Add("filter", firmwareUpdateApiFilter("eq", "product", unifiOSServerProduct))
	q.Add("filter", firmwareUpdateApiFilter("eq", "platform", linuxX64Platform))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var respData firmwareUpdateApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, nil, err
	}

	for _, fw := range respData.Embedded.Firmware {
		if fw.Platform != linuxX64Platform {
			continue
		}
		return fw.Version.Core(), fw.Links.Data.Href, nil
	}

	return nil, nil, fmt.Errorf("no unifi-os-server firmware found for platform %q", linuxX64Platform)
}
