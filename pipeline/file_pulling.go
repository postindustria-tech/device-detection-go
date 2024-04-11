package pipeline

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// scheduleFilePulling schedules the file pulling
func (p *Pipeline) scheduleFilePulling() {
	// recover from panic
	// if panic occurs, we will log the error and restart the file pulling
	defer func() {
		if r := recover(); r != nil {
			p.logger.Printf("error occurred when pulling data: %v", r)
			go p.scheduleFilePulling()
		}
	}()

	nextIterationInMs := p.dataFilePullEveryMs

	isFirstRun := true

	for {
		// if this is the first run, we don't need to wait
		if isFirstRun {
			isFirstRun = false
		} else {
			<-time.After(time.Duration(nextIterationInMs) * time.Millisecond)
		}

		p.logger.Printf("Pulling data from %s", p.dataFileUrl)

		fileResponse, err := doDataFileRequest(p.dataFileUrl)
		if err != nil {
			p.logger.Printf("failed to pull data file: %v", err)
			if fileResponse != nil && fileResponse.retryAfter > 0 {
				p.logger.Printf("received retry-after, retrying after %d seconds", fileResponse.retryAfter)

				// retry after the specified time
				nextIterationInMs = fileResponse.retryAfter * 1000
			} else {
				p.logger.Printf("retrying after 1 second")
				// retry after 1 second, since we have unhandled error
				// this can happen from network stutter or something else
				nextIterationInMs = 1000
			}
			continue
		}

		p.logger.Printf("data file pulled successfully")

		// write to dataFile
		err = os.WriteFile(p.dataFile, fileResponse.buffer.Bytes(), 0644)
		if err != nil {
			p.logger.Printf("failed to write data file: %v", err)
			// retry after 1 second, since we have unhandled error
			// this can happen from disk write error or something else
			nextIterationInMs = 1000
			continue
		}
		p.logger.Printf("data file written successfully: %d bytes", fileResponse.buffer.Len())

		err = p.manager.ReloadFromOriginalFile()
		p.logger.Printf("data file reloaded successfully")
		if err != nil {
			p.logger.Printf("failed to reload data file: %v", err)
			// retry after 1 second, since we have unhandled error
			// this can happen from reload error or something else
			nextIterationInMs = 1000
			continue
		}
		// reset nextIterationInMs
		nextIterationInMs = p.dataFilePullEveryMs
	}
}

type FileResponse struct {
	buffer     bytes.Buffer
	retryAfter int
}

// doDataFileRequest performs the data file request
func doDataFileRequest(url string) (*FileResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//handle retry after
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		retryAfterSeconds, err := strconv.Atoi(retryAfter)
		if err != nil || len(retryAfter) == 0 {
			return nil, fmt.Errorf("failed to parse Retry-After header: %v", err)
		}

		return &FileResponse{
			retryAfter: retryAfterSeconds,
		}, fmt.Errorf("received 429, retrying after %d seconds", retryAfterSeconds)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to pull data file: %s", resp.Status)
	}

	contentMD5 := resp.Header.Get("Content-MD5")

	fmt.Printf("Content-MD5: %s\n", contentMD5)

	isValid, err := validateMd5(resp.Body, contentMD5)
	if err != nil {
		return nil, fmt.Errorf("failed to validate MD5: %v", err)
	}

	if !isValid {
		return nil, fmt.Errorf("MD5 validation failed")
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return &FileResponse{
		buffer:     buffer,
		retryAfter: 0,
	}, nil
}

func validateMd5(body io.ReadCloser, hash string) (bool, error) {
	h := md5.New()
	_, err := io.Copy(h, body)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString(h.Sum(nil)) == hash, nil
}
