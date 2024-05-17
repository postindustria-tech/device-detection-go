package onpremise

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/51Degrees/device-detection-go/v4/dd"
)

func unzipAndSaveToTempFile(name string) (*os.File, error) {
	mockHashMutex.Lock()
	file, err := os.Open("mock_hash.gz")
	defer mockHashMutex.Unlock()
	defer file.Close()
	if err != nil {
		return nil, err
	}
	gReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gReader.Close()

	uncompressed, err := io.ReadAll(gReader)
	err = os.WriteFile(name, uncompressed, 0644)
	if err != nil {
		return nil, err
	}

	uFile, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer uFile.Close()

	return uFile, nil
}

func TestExternalFileChangedReplace(t *testing.T) {
	config := dd.NewConfigHash(dd.Balanced)

	tempFile, err := unzipAndSaveToTempFile("TestExternalFileChangedReplace.hash")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	engine, err := New(
		config,
		WithDataFile("TestExternalFileChangedReplace.hash"),
		WithFileWatch(true),
		WithAutoUpdate(false),
	)
	if err != nil {
		t.Fatalf("Error creating engine: %v", err)
	}
	defer engine.Stop()
	defer func() {
		files, err := filepath.Glob("*-TestExternalFileChangedReplace.hash")
		if err != nil {
			t.Fatalf("Error listing files: %v", err)
		}

		for _, f := range files {
			err = os.Remove(f)
			if err != nil {
				t.Fatalf("Error removing file: %v", err)
			}
		}
	}()

	mockEvidence := []Evidence{
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Arch",
			Value:  "x86",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Model",
			Value:  "Intel",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Mobile",
			Value:  "?0",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Platform",
			Value:  "Windows",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Platform-Version",
			Value:  "10.0",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Full-Version-List",
			Value:  "58.0.3029.110",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua",
			Value:  `"\"Chromium\";v=\"91.0.4472.124\";a=\"x86\";p=\"Windows\";rv=\"91.0\""`,
		},
	}

	resultsHash, err := engine.Process(mockEvidence)
	defer resultsHash.Free()
	if err != nil {
		t.Fatalf("Failed to process evidence: %v", err)
	}

	browser, err := resultsHash.ValuesString("BrowserName", ",")
	if err != nil {
		log.Fatalf("Failed to get BrowserName: %v", err)
	}

	if browser != "Chromium Project" {
		t.Fatalf("Expected BrowserName to be Chromium Project, got %s", browser)
	}
	tempFile2, err := unzipAndSaveToTempFile("TestExternalFileChangedReplace.hash")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer tempFile2.Close()
	defer os.Remove(tempFile2.Name())
	resultsHash2, err := engine.Process(mockEvidence)
	defer resultsHash2.Free()
	if err != nil {
		t.Fatalf("Failed to process evidence: %v", err)
	}

	browser, err = resultsHash2.ValuesString("BrowserName", ",")
	if err != nil {
		log.Fatalf("Failed to get BrowserName: %v", err)
	}

	if browser != "Chromium Project" {
		t.Fatalf("Expected BrowserName to be Chromium Project, got %s", browser)
	}

	<-time.After(4 * time.Second)
	if engine.fileExternallyChangedCount != 1 {
		t.Fatalf("Expected 1 file change, got %d", engine.fileExternallyChangedCount)
	}
}

func TestExternalFileChangedMv(t *testing.T) {
	config := dd.NewConfigHash(dd.Balanced)

	tempFile, err := unzipAndSaveToTempFile("TestExternalFileChangedMv.hash")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	engine, err := New(
		config,
		WithDataFile("TestExternalFileChangedMv.hash"),
		WithFileWatch(true),
		WithAutoUpdate(false),
	)
	if err != nil {
		t.Fatalf("Error creating engine: %v", err)
	}
	defer engine.Stop()
	defer func() {
		files, err := filepath.Glob("*-TestExternalFileChangedMv.hash")
		if err != nil {
			t.Fatalf("Error listing files: %v", err)
		}

		for _, f := range files {
			err = os.Remove(f)
			if err != nil {
				t.Fatalf("Error removing file: %v", err)
			}
		}
	}()
	tempFile2, err := unzipAndSaveToTempFile("TestExternalFileChangedMv2.hash")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile2.Name())
	defer tempFile2.Close()
	err = os.Rename(tempFile2.Name(), tempFile.Name())
	if err != nil {
		t.Fatalf("Error renaming file: %v", err)
	}

	mockEvidence := []Evidence{
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Arch",
			Value:  "x86",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Model",
			Value:  "Intel",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Mobile",
			Value:  "?0",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Platform",
			Value:  "Windows",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Platform-Version",
			Value:  "10.0",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua-Full-Version-List",
			Value:  "58.0.3029.110",
		},
		{
			Prefix: dd.HttpHeaderString,
			Key:    "Sec-Ch-Ua",
			Value:  `"\"Chromium\";v=\"91.0.4472.124\";a=\"x86\";p=\"Windows\";rv=\"91.0\""`,
		},
	}

	<-time.After(5 * time.Second)
	resultsHash, err := engine.Process(mockEvidence)
	defer resultsHash.Free()
	if err != nil {
		t.Fatalf("Failed to process evidence: %v", err)
	}

	browser, err := resultsHash.ValuesString("BrowserName", ",")
	if err != nil {
		log.Fatalf("Failed to get BrowserName: %v", err)
	}

	if browser != "Chromium Project" {
		t.Fatalf("Expected BrowserName to be Chromium Project, got %s", browser)
	}

	<-time.After(4 * time.Second)
	if engine.fileExternallyChangedCount != 1 {
		t.Fatalf("Expected 1 file change, got %d", engine.fileExternallyChangedCount)
	}
}
