package onpremise

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const retryMs = 1000

var (
	ErrTooManyRetries = errors.New("too many retries to pull data file")
)

// scheduleFilePulling schedules the file pulling
func (p *Engine) scheduleFilePulling() {
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
	retryAttempts := 0

	for {
		select {
		case <-p.stopCh:
			return
		default:
			if retryAttempts > 5 && !p.isManagerInitialized {
				p.rdySignal <- ErrTooManyRetries
				return
			}

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
					if p.isManagerInitialized && p.fileSynced {
						p.rdySignal <- nil
					}
					p.logger.Printf("received retry-after, retrying after %d seconds", fileResponse.retryAfter)
					// retry after the specified time
					nextIterationInMs = fileResponse.retryAfter * 1000
				} else {
					p.logger.Printf("retrying after 1 second")
					retryAttempts += 1
					// retry after 1 second, since we have unhandled error
					// this can happen from network stutter or something else
					nextIterationInMs = retryMs
				}
				continue
			}

			p.logger.Printf("data file pulled successfully: %d bytes", fileResponse.buffer.Len())

			err = p.createDatafileIfNotExists()
			if err != nil {
				p.logger.Printf("failed to create data file: %v", err)
				// retry after 1 second, since we have unhandled error
				// this can happen from disk write error or something else
				retryAttempts += 1
				nextIterationInMs = retryMs
				continue
			}

			// write to dataFile
			err = os.WriteFile(p.dataFile, fileResponse.buffer.Bytes(), os.ModeAppend)
			if err != nil {
				p.logger.Printf("failed to write data file: %v", err)
				// retry after 1 second, since we have unhandled error
				// this can happen from disk write error or something else
				retryAttempts += 1
				nextIterationInMs = retryMs
				continue
			}
			p.logger.Printf("data file written successfully: %d bytes", fileResponse.buffer.Len())

			if !p.isManagerInitialized {
				err = p.initializeManager()
				if err != nil {
					p.logger.Printf("failed to initialize manager: %v", err)
					// retry after 1 second, since we have unhandled error
					// this can happen from manager initialization error or something else
					retryAttempts += 1
					nextIterationInMs = retryMs
					continue
				}
			}

			err = p.manager.ReloadFromOriginalFile()
			if err != nil {
				p.logger.Printf("failed to reload data file: %v", err)
				// retry after 1 second, since we have unhandled error
				// this can happen from reload error or something else
				retryAttempts += 1
				nextIterationInMs = retryMs
				continue
			}
			p.logger.Printf("data file reloaded successfully")

			if !p.fileSynced {
				p.fileSynced = true
				p.rdySignal <- nil
			}

			// reset nextIterationInMs
			nextIterationInMs = p.dataFilePullEveryMs
			p.totalFilePulls += 1
			retryAttempts = 0
		}
	}
}

type FileResponse struct {
	buffer     *bytes.Buffer
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

	buffer := bytes.NewBuffer(make([]byte, 0))
	_, err = buffer.ReadFrom(resp.Body)
	contentMD5 := resp.Header.Get("Content-MD5")

	fileBytes, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress file: %v", err)
	}

	uncompressedBuffer := bytes.NewBuffer(make([]byte, 0))
	_, err = uncompressedBuffer.ReadFrom(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed file: %v", err)
	}
	uncompressedBytes := uncompressedBuffer.Bytes()

	isValid, err := validateMd5(uncompressedBytes, contentMD5)
	if err != nil {
		return nil, fmt.Errorf("failed to validate MD5: %v", err)
	}

	if !isValid {
		return nil, fmt.Errorf("MD5 validation failed")
	}

	return &FileResponse{
		buffer:     bytes.NewBuffer(uncompressedBytes),
		retryAfter: 0,
	}, nil
}

func validateMd5(val []byte, hash string) (bool, error) {
	h := md5.New()
	_, err := h.Write(val)
	if err != nil {
		return false, err
	}

	strVal := hex.EncodeToString(h.Sum(nil))

	return strVal == hash, nil
}
