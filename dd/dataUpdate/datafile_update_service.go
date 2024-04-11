package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type DataFileUpdateService struct {
	Updating bool
}

func (service *DataFileUpdateService) UpdateDataFile(df DataFile) error {
	if service.Updating {
		slog.Warn("Updating data is already in progress")
		return nil
	}
	service.Updating = true
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered: ", r)
		}
		service.Updating = false
		service.CheckNextUpdate(df)
	}()

	req, err := http.NewRequest("GET", df.urlFormatter(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	if df.VerifyIfModifiedSince && df.Path != "" {
		publishDate, err := df.getPublishedDate()
		if err == nil {
			req.Header.Set("If-Modified-Since", publishDate.Format(http.TimeFormat))
		} else {
			slog.Warn("Couldn't verify modified since")
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			return fmt.Errorf("too many requests to %s", df.urlFormatter())
		case http.StatusNotModified:
			return fmt.Errorf("no data update available at %s", df.urlFormatter())
		case http.StatusForbidden:
			return fmt.Errorf("access denied to %s", df.urlFormatter())
		default:
			return fmt.Errorf("response error from %s: %s", df.urlFormatter(), resp.Status)
		}
	}

	contentMD5 := resp.Header.Get("Content-MD5")
	if contentMD5 != "" {
		verify, err := service.VerifyMD5(resp.Body, contentMD5)
		if err != nil {
			return fmt.Errorf("error verifying MD5: %w", err)
		}
		if !verify {
			errMsg := fmt.Sprintf("response header MD5 '%s' does not match file's", contentMD5)
			return errors.New(errMsg)
		}
	}

	err = service.SaveDataFile(df, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving datafile: %w", err)
	}
	return nil
}

func (service *DataFileUpdateService) SaveDataFile(df DataFile, data io.ReadCloser) error {
	err := os.MkdirAll(df.TempPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating temp folder: %w", err)
	}

	tempFilePath := df.TempPath + "/" + df.Identifier
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer tempFile.Close()

	_, err = tempFile.ReadFrom(data)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	err = os.Rename(tempFilePath, df.Path+"/"+df.Identifier)
	if err != nil {
		return fmt.Errorf("error moving file: %w", err)
	}

	return nil
}

func (service *DataFileUpdateService) CheckNextUpdate(df DataFile) error {
	if df.AutoUpdate {
		randomOffset := rand.Intn(abs(df.UpdateTimeMaximumRandomisation) + 1)
		interval := 60*abs(df.PollingInterval) + randomOffset

		time.AfterFunc(time.Second*time.Duration(interval), func() { service.UpdateDataFile(df) })
		//TODO error return
	}
	return nil
}

func (service *DataFileUpdateService) RegisterDataFile(df DataFile) error {
	if df.UpdateOnStart {
		return service.UpdateDataFile(df)
	} else {
		return service.CheckNextUpdate(df)
	}
}

func (service *DataFileUpdateService) VerifyMD5(data io.ReadCloser, md5Hash string) (bool, error) {
	h := md5.New()
	if _, err := io.Copy(h, data); err != nil {
		return false, fmt.Errorf("error copying file: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)) == md5Hash, nil
}

func abs(x int) int {
	if x == math.MinInt64 {
		return math.MaxInt64
	} else if x < 0 {
		return -x
	}
	return x
}
