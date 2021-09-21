package dd

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// test User-Agents
const uaMobile = "Mozilla/5.0 (iPhone; CPU iPhone OS 7_1 like Mac OS X) " +
	"AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D167 " +
	"Safari/9537.53"

// User-Agent string of Firefox Web browser version 41 on desktop.
const uaDesktop = "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) " +
	"Gecko/20100101 Firefox/41.0"

// User-Agent string of a MediaHub device.
const uaMediaHub = "Mozilla/5.0 (Linux; Android 4.4.2; X7 Quad Core " +
	"Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 " +
	"Chrome/30.0.0.0 Safari/537.36"

// File constants
const testLiteDataFile = "51Degrees-LiteV4.1.hash"
const testEnterpriseDataFile = "Enterprise-HashV41.hash"

// Test variables
var testDataFilePath string = ""

// getFilePath take a directory and search for the target file in that directory,
// including the subdirectories
func getFilePath(dir string, names []string) (found string, err error) {
	path, e := os.Getwd()
	if e != nil {
		log.Fatalln("Failed to get current directory.")
	}

	// Target directory
	path += "/" + dir

	// walk through the directory
	foundPath := ""
	e = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error \"%s\" occured on path \"%s\".\n",
				err,
				path)
			return err
		}
		for _, name := range names {
			if strings.EqualFold(info.Name(), name) {
				foundPath = path
				// Return error Exist here to stop the walk
				return fs.ErrExist
			}
		}
		return nil
	})

	// Return error not exist to indicate no matching find found.
	if foundPath == "" && e != fs.ErrExist {
		e = fs.ErrNotExist
	}
	return foundPath, e
}

func TestMain(m *testing.M) {
	// Search for test data files
	dataFiles := []string{testLiteDataFile, testEnterpriseDataFile}
	filePath, err := getFilePath(
		"device-detection-cxx",
		dataFiles)
	if err != nil && err != fs.ErrExist {
		fileList := ""
		for _, file := range dataFiles {
			if fileList != "" {
				fileList += ", "
			}
			fileList += file
		}
		log.Fatalf("Could not find any file that matches any of \"%s\".\n",
			fileList)
	}

	testDataFilePath = filePath
	os.Exit(m.Run())
}
