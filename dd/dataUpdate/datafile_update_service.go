package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type DataFileUpdateService struct {
	Updating bool
}

func (service *DataFileUpdateService) UpdateDataFile(df DataFile) {
	if service.Updating {
		//TODO log message
		return
	}
	service.Updating = true
	defer func() {
		service.Updating = false
		service.CheckNextUpdate(df)
	}()

	req, err := http.NewRequest("GET", df.urlFormatter(), nil)
	if err != nil {
		//TODO: err
		log.Fatal(err)
	}
	if df.VerifyIfModifiedSince {
		req.Header.Set("If-Modified-Since", df.DataPublishedDateTime.Format(http.TimeFormat))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: err
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
		case http.StatusNotModified:
		case http.StatusForbidden:
		default:
		}
		//TODO: err
		return
	}

	contentMD5 := resp.Header.Get("Content-MD5")
	if contentMD5 != "" && !service.VerifyMD5(resp.Body, contentMD5) {
		//TODO: err if VerifyMD5 false
		return
	}

	service.SaveDataFile(df, resp.Body)
}

func (service *DataFileUpdateService) SaveDataFile(df DataFile, data io.ReadCloser) {
	err := os.MkdirAll(df.TempPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		//TODO: err
	}

	tempFilePath := df.TempPath + "/" + df.Identifier
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Fatal(err)
		//TODO: err
	}
	defer tempFile.Close()

	_, err = tempFile.ReadFrom(data)
	if err != nil {
		log.Fatal(err)
		//TODO: err
	}

	err = os.Rename(tempFilePath, df.Path+"/"+df.Identifier)
	if err != nil {
		log.Fatal(err)
		//TODO: err
	}
}

func (service *DataFileUpdateService) CheckNextUpdate(df DataFile) {
	if df.AutoUpdate {
		randomOffset := rand.Intn(abs(df.UpdateTimeMaximumRandomisation) + 1)
		interval := 60*abs(df.PollingInterval) + randomOffset

		time.AfterFunc(time.Second*time.Duration(interval), func() { service.UpdateDataFile(df) })
	}
}

func (service *DataFileUpdateService) RegisterDataFile(df DataFile) {
	if df.UpdateOnStart {
		service.UpdateDataFile(df)
	} else {
		service.CheckNextUpdate(df)
	}
}

func (service *DataFileUpdateService) VerifyMD5(data io.ReadCloser, md5Hash string) bool {
	h := md5.New()
	if _, err := io.Copy(h, data); err != nil {
		//TODO: err
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil)) == md5Hash
}

func abs(x int) int {
	if x == math.MinInt64 {
		return math.MaxInt64
	} else if x < 0 {
		return -x
	}
	return x
}
