package dd

// #include <string.h>
// #include "./device-detection-cxx/src/common-cxx/fiftyone.h"
import "C"
import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetFilePath take a directory and search for the target file in that directory,
// including the subdirectories. The directory is a relative path to the current
// directory of the caller file.
func GetFilePath(dir string, names []string) (found string, err error) {
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
