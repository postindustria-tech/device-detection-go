package onpremise

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/51Degrees/device-detection-go/v4/dd"
)

func TestCustomProvider(t *testing.T) {
	mockServer := newMockDataFileServer(10 * time.Second)
	defer mockServer.Close()

	cases := []struct {
		name          string
		engineOptions []EngineOptions
		expectedError string
	}{
		{
			name: "with license key and custom url",
			engineOptions: []EngineOptions{
				WithLicenseKey("123"),
				WithDataUpdateUrl(mockServer.URL + "/datafile"),
				WithPollingInterval(2),
				WithFileWatch(false),
			},
			expectedError: "",
		},
		{
			name: "with product and custom url",
			engineOptions: []EngineOptions{
				WithProduct("MyProduct"),
				WithDataUpdateUrl(mockServer.URL + "/datafile"),
				WithPollingInterval(2),
				WithFileWatch(false),
			},
			expectedError: "",
		},
		{
			name: "Invalid url",
			engineOptions: []EngineOptions{
				WithDataUpdateUrl("dsoahdsakjhd"),
				WithPollingInterval(2),
				WithFileWatch(false),
			},
			expectedError: `parse "dsoahdsakjhd": invalid URI for request`,
		},
		{
			name: "with custom url",
			engineOptions: []EngineOptions{
				WithDataUpdateUrl(mockServer.URL + "/datafile"),
				WithPollingInterval(2),
				WithFileWatch(false),
			},
			expectedError: "",
		},
		{
			name: "no data file",
			engineOptions: []EngineOptions{
				WithDataUpdateUrl(mockServer.URL + "/datafile"),
				WithPollingInterval(2),
				WithFileWatch(false),
			},
			expectedError: ErrNoDataFileProvided.Error(),
		},
	}

	for _, c := range cases {
		tempFile, err := unzipAndSaveToTempFile(fmt.Sprintf("test_%s.hash", c.name))
		if err != nil {
			t.Fatalf("Error creating temp file: %v", err)
		}

		options := append(c.engineOptions, WithDataFile(tempFile.Name()))

		defer tempFile.Close()
		defer os.Remove(tempFile.Name())
		config := dd.NewConfigHash(dd.Balanced)
		engine, err := New(
			config,
			options...,
		)
		if engine != nil {
			engine.Stop()
		}
		<-time.After(3 * time.Second)
		if err != nil && err.Error() != c.expectedError {
			t.Errorf("Failed case %s: %v", c.name, err)
		}
	}
}

func testProperties(t *testing.T, properties []string, values []string, noValueProps []string) {
	config := dd.NewConfigHash(dd.Balanced)
	engine, err := New(
		config,
		WithDataFile("../51Degrees-LiteV4.1.hash"),
		WithAutoUpdate(false),
		WithProperties(properties),
	)

	if err != nil {
		t.Fatalf("Error creating engine: %v", err)
	}
	result, err := engine.Process([]Evidence{
		{
			dd.HttpHeaderString,
			"User-Agent",
			"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.72 Mobile Safari/537.36"},
	})
	if err != nil {
		t.Errorf("Error processing request: %v", err)
	}
	defer result.Free()
	for i := 0; i < len(properties); i++ {
		value, _ := result.ValuesString(properties[i], "")
		if value != values[i] {
			t.Errorf("Error processing request: expected %s, got %s", values[i], value)
		}
	}

	for i := 0; i < len(noValueProps); i++ {
		value, _ := result.ValuesString(noValueProps[i], "")
		if value != "" {
			t.Errorf("Error processing request: expected no value, got %s", value)
		}
	}
}

func TestWithPropertiesNoPlatform(t *testing.T) {
	testProperties(t, []string{"IsMobile", "BrowserName"}, []string{"True", "Chrome Mobile"}, []string{"PlatformName"})
}

func TestWithPropertiesPlatform(t *testing.T) {
	testProperties(t, []string{"PlatformName"}, []string{"Android"}, []string{"IsMobile", "BrowserName"})
}

func TestPropertiesDefault(t *testing.T) {
	testProperties(t, nil, []string{"True", "Chrome Mobile", "Android"}, nil)
}
