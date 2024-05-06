package onpremise

import (
	"fmt"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"os"
	"testing"
	"time"
)

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
				WithLicenceKey("123"),
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
