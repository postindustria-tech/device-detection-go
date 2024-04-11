package main

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

type DataFile struct {
	Identifier                     string    // Name of the datafile
	TempPath                       string    // Temporary datafile path
	Path                           string    // Path to the datafile
	AutoUpdate                     bool      // Whether to automatically update the datafile
	PollingInterval                int       // How often to poll for updates to the datafile (minutes)
	UpdateTimeMaximumRandomisation int       // Maximum randomisation offset in seconds to polling time interval
	UpdateOnStart                  bool      // Whether to update the datafile as soon as it is initialised
	UpdateURLParams                URLParams // URL params for the update
	DataPublishedDateTime          time.Time
	VerifyIfModifiedSince          bool
	LastUpdateTryTime              time.Time
	Updating                       bool
}

//TODO constructor for DataFile?

type URLParams struct {
	url         string // URL path to  51Degrees Distributor service
	LicenseKeys string /*
		Pipe-separated list of 51Degrees License Keys.
		You must have a valid License Key for the data file that you are trying to download.
	*/
	product string /*
	  The 51Degrees 'product' that the data file is part of.
	   Accepts "Premium", "Enterprise", "Lite", "TAC", "V4Free", "V4Enterprise", "V4TAC" values.
	*/
	urlType string /*
	   The type of data file to download.
	   Accepts "BinaryV3", "BinaryV3Uncompressed", "BinaryV32", "BinaryV32Uncompressed",
	   "BinaryV32UAT", "TrieV32", "HashTrieV34", "CSV", "HashV41" values.
	*/
}

func NewURLParams(urlString, LicenseKeys, product, urlType string) (URLParams, error) {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return URLParams{}, fmt.Errorf("invalid url: %w", err)
	}
	switch product {
	case "Premium", "Enterprise", "Lite", "TAC", "V4Free", "V4Enterprise", "V4TAC":
	default:
		return URLParams{}, fmt.Errorf("invalid product value")
	}
	switch urlType {
	case "BinaryV3", "BinaryV3Uncompressed", "BinaryV32", "BinaryV32Uncompressed", "BinaryV32UAT", "TrieV32", "HashTrieV34", "CSV", "HashV41":
	default:
		return URLParams{}, fmt.Errorf("invalid type value")

	}
	return URLParams{url: urlString, LicenseKeys: LicenseKeys, product: product, urlType: urlType}, nil
}

func (param *URLParams) ToString() string {
	u, _ := url.Parse(param.url)

	query := u.Query()
	query.Set("LicenseKeys", param.LicenseKeys)
	query.Set("Product", param.product)
	query.Set("Type", param.urlType)

	u.RawQuery = query.Encode()
	return u.String()
}

func (df *DataFile) urlFormatter() string {
	return df.UpdateURLParams.ToString()
}

func (df *DataFile) getPublishedDate() (time.Time, error) {
	var publishedDate = df.DataPublishedDateTime
	if publishedDate.IsZero() {
		if df.Path == "" {
			return time.Time{}, fmt.Errorf("no publish date exists")
		}
		fileInfo, err := os.Stat(df.Path)
		if err != nil {
			return time.Time{}, fmt.Errorf("can't get publish date: %w", err)
		}
		publishedDate = fileInfo.ModTime()
	}
	return publishedDate, nil
}
