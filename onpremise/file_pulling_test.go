package onpremise

import (
	"bytes"
	"errors"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"log"

	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func newMockDataFileServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Open the file for reading
				file, err := os.Open("./mock_hash.gz")
				if err != nil {
					log.Printf("Failed to open file: %v", err)
					return
				}
				defer file.Close()

				buffer := bytes.NewBuffer(make([]byte, 0))
				_, err = io.Copy(buffer, file)
				if err != nil {
					log.Printf("Failed to read file: %v", err)
					return
				}

				w.Header().Add("Content-MD5", "daebfa89ddefac8e6c4325c38f129504")
				w.Header().Add("Content-Length", strconv.Itoa(buffer.Len()))

				w.WriteHeader(http.StatusOK)

				// Write the file to response
				_, err = io.Copy(w, buffer)
				if err != nil {
					log.Printf("Failed to write file: %v", err)
					return
				}

			},
		),
	)
}

func TestFilePulling(t *testing.T) {
	server := newMockDataFileServer()
	defer server.Close()

	config := dd.NewConfigHash(dd.Balanced)

	pl, err := New(
		config,
		WithDataUpdateUrl(
			server.URL+"/datafile",
			2,
		),
		SetMaxRetries(2),
	)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	<-time.After(5000 * time.Millisecond)

	if pl.totalFilePulls != 2 {
		t.Fatalf("Expected 2 file pulls, got %d", pl.totalFilePulls)
	}

	resultsHash, err := pl.Process(
		[]Evidence{
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
		},
	)
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

	log.Printf("BrowserName: %s", browser)

	deviceType, err := resultsHash.ValuesString("DeviceType", ",")
	if err != nil {
		log.Fatalf("Failed to get DeviceType: %v", err)
	}

	if deviceType != "Desktop" {
		t.Fatalf("Expected DeviceType to be Desktop, got %s", deviceType)
	}
}

func TestTooManyRetries(t *testing.T) {
	config := dd.NewConfigHash(dd.Balanced)
	_, err := New(
		config,
		SetLicenceKey("123"),
		SetProduct("MyProduct"),
		SetMaxRetries(5),
	)
	if !errors.Is(err, ErrTooManyRetries) {
		t.Fatalf("Expected error to be ErrTooManyRetries, got %v", err)
	}
}

func newMockServerModifiedSince() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				val := r.Header.Get("If-Modified-Since")
				if val == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusNotModified)
			},
		),
	)
}

func TestIfModifiedSince(t *testing.T) {
	server := newMockServerModifiedSince()

	defer server.Close()

	timeStamp := time.Now()
	_, err := doDataFileRequest(
		server.URL, logWrapper{
			enabled: true,
			logger:  DefaultLogger,
		}, &timeStamp,
	)

	if !errors.Is(err, ErrFileNotModified) {
		t.Errorf("Expected error to be ErrFileNotModified, got %v", err)
	}
}
