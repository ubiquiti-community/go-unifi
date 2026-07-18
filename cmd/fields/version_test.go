package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLatestUnifiVersion(t *testing.T) {
	assert := assert.New(t)
	_ = require.New(t)

	fwVersion, err := version.NewVersion("7.3.83+atag-7.3.83-19645")
	assert.NoError(err)

	fwDownload, err := url.Parse(
		"https://fw-download.ubnt.com/data/unifi-controller/c31c-debian-7.3.83-c9249c913b91416693b869b9548850c3.deb",
	)
	assert.NoError(err)

	respData := firmwareUpdateApiResponse{
		Embedded: firmwareUpdateApiResponseEmbedded{
			Firmware: []firmwareUpdateApiResponseEmbeddedFirmware{
				{
					Channel:  releaseChannel,
					Created:  "2023-02-06T08:55:31+00:00",
					Id:       "c9249c91-3b91-4166-93b8-69b9548850c3",
					Platform: debianPlatform,
					Product:  unifiControllerProduct,
					Version:  fwVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{
							Href: fwDownload,
						},
					},
				},
				{
					Channel:  releaseChannel,
					Created:  "2023-02-06T08:51:36+00:00",
					Id:       "2a600108-7f79-4b3e-b6e0-4dd262460457",
					Platform: "document",
					Product:  unifiControllerProduct,
					Version:  fwVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{
							Href: nil,
						},
					},
				},
				{
					Channel:  releaseChannel,
					Created:  "2023-02-06T08:51:37+00:00",
					Id:       "9d2d413d-36ce-4742-a10d-4351aac6f08d",
					Platform: "windows",
					Product:  unifiControllerProduct,
					Version:  fwVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{
							Href: nil,
						},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		assert.Contains(query["filter"], firmwareUpdateApiFilter("eq", "channel", releaseChannel))
		assert.Contains(
			query["filter"],
			firmwareUpdateApiFilter("eq", "product", unifiControllerProduct),
		)
		assert.Contains(query["filter"], firmwareUpdateApiFilter("lt", "version", maxVersion))

		resp, err := json.Marshal(respData)
		assert.NoError(err)

		_, err = rw.Write(resp)
		assert.NoError(err)
	}))
	defer server.Close()

	firmwareUpdateApi = server.URL
	gotVersion, gotDownload, err := latestUnifiVersion()
	assert.NoError(err)

	assert.Equal(fwVersion.Core(), gotVersion)
	assert.Equal(fwDownload, gotDownload)
}

func TestLatestUnifiOSVersion(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	osVersion, err := version.NewVersion("v5.1.21")
	require.NoError(err)

	osDownload, err := url.Parse(
		"https://fw-download.ubnt.com/data/unifi-os-server/f5e2-linux-x64-5.1.21-a400c9c6-8328-4634-b223-ebfcf742720a.21-x64",
	)
	require.NoError(err)

	respData := firmwareUpdateApiResponse{
		Embedded: firmwareUpdateApiResponseEmbedded{
			Firmware: []firmwareUpdateApiResponseEmbeddedFirmware{
				{
					Channel:  releaseChannel,
					Platform: linuxX64Platform,
					Product:  unifiOSServerProduct,
					Version:  osVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{
							Href: osDownload,
						},
					},
				},
				{
					Channel:  releaseChannel,
					Platform: "macOS-dmg-amd64",
					Product:  unifiOSServerProduct,
					Version:  osVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{Href: nil},
					},
				},
				{
					Channel:  releaseChannel,
					Platform: linuxX64Platform,
					Product:  "unifi-controller", // wrong product
					Version:  osVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{Href: nil},
					},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		assert.Contains(query["filter"], firmwareUpdateApiFilter("eq", "channel", releaseChannel))
		assert.Contains(query["filter"], firmwareUpdateApiFilter("eq", "product", unifiOSServerProduct))
		assert.Contains(query["filter"], firmwareUpdateApiFilter("eq", "platform", linuxX64Platform))

		resp, err := json.Marshal(respData)
		assert.NoError(err)
		_, err = rw.Write(resp)
		assert.NoError(err)
	}))
	defer srv.Close()

	firmwareUpdateApi = srv.URL
	gotVersion, gotDownload, err := latestUnifiOSVersion()
	require.NoError(err)

	assert.Equal(osVersion.Core(), gotVersion)
	assert.Equal(osDownload, gotDownload)
}
