package onpremise

import (
	"github.com/51Degrees/device-detection-go/v4/dd"
	"testing"
)

func TestDefaultProvider(t *testing.T) {
	config := dd.NewConfigHash(dd.Balanced)
	engine, err := New(
		config,
		SetLicenceKey("123"),
		SetProduct("MyProduct"),
	)
	if err != nil {
		t.Errorf("Failed to create engine: %v", err)
	}

	if engine.dataFileUrl != "https://distributor.51degrees.com/api/v2/download?Download=True&LicenseKeys=123&Product=MyProduct&Type=HashV41" {
		t.Errorf("Failed to set default data file url %s", engine.dataFileUrl)
	}

	engine.Stop()
}

func TestCustomProvider(t *testing.T) {
	mockServer := newMockDataFileServer()
	defer mockServer.Close()

	cases := []struct {
		name          string
		engineOptions []EngineOptions
		expectedError string
	}{
		{
			name: "with licence key and custom url",
			engineOptions: []EngineOptions{
				SetLicenceKey("123"),
				WithDataUpdateUrl(mockServer.URL+"/datafile", 2000),
			},
			expectedError: "",
		},
		{
			name: "with product and custom url",
			engineOptions: []EngineOptions{
				SetProduct("MyProduct"),
				WithDataUpdateUrl(mockServer.URL+"/datafile", 2000),
			},
			expectedError: "",
		},
		{
			name: "Invalid url",
			engineOptions: []EngineOptions{
				WithDataUpdateUrl("dsoahdsakjhd", 2000),
			},
			expectedError: `parse "dsoahdsakjhd": invalid URI for request`,
		},
		{
			name: "with custom url",
			engineOptions: []EngineOptions{
				WithDataUpdateUrl(mockServer.URL+"/datafile", 2000),
			},
			expectedError: "",
		},
		{
			name: "with data file",
			engineOptions: []EngineOptions{
				WithDataFile("non existing file"),
			},
			expectedError: "failed to get file path: file does not exist",
		},
	}

	config := dd.NewConfigHash(dd.Balanced)

	for _, c := range cases {
		engine, err := New(
			config,
			c.engineOptions...,
		)
		if engine != nil {
			engine.Stop()
		}
		if err != nil && err.Error() != c.expectedError {
			t.Errorf("Failed case %s: %v", c.name, err)
		}

	}
}
