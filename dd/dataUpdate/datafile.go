package main

import (
	"time"
)

type DataFile struct {
	Identifier                     string // Name of the datafile
	TempPath                       string // Temporary datafile path
	Path                           string // Path to the datafile
	AutoUpdate                     bool   // Whether to automatically update the datafile
	PollingInterval                int    // How often to poll for updates to the datafile (minutes)
	UpdateTimeMaximumRandomisation int    // Maximum randomisation offset in seconds to polling time interval
	UpdateOnStart                  bool   // Whether to update the datafile as soon as it is initialised
	DataUpdateUrl                  string // TODO: Should be object?
	DataPublishedDateTime          time.Time
	VerifyIfModifiedSince          bool
	LastUpdateTryTime              time.Time
	Updating                       bool
}

func (df *DataFile) urlFormatter() string {
	// by default just return the url
	return df.DataUpdateUrl
}

// func (df *DataFile) getPublishedDate() time.Time {
// 	var publishedDate = df.DataPublishedDateTime
// 	if publishedDate.IsZero() {
// 		fileInfo, err := os.Stat(df.Path)
// 		if err != nil {

// 		}
// 		publishedDate = fileInfo.ModTime()
// 	}
// 	return publishedDate
// }
