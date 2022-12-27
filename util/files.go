package util

import (
	"flag"
	"fmt"
	"github.com/its-my-data/doubak/proto"
	"github.com/mengzhuo/cookiestxt"
	"html"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const CollectorPathPrefix = "collector/"
const ItemPathPrefix = "items/"

// GetPathWithCreation returns the concatenated path with output path and have them created in advance.
func GetPathWithCreation(subdirs string) (string, error) {
	baseDir := flag.Lookup(proto.Flag_output_dir.String()).Value.String()
	return GetPathWithCreationWithBase(baseDir, subdirs)
}

// GetPathWithCreationWithBase returns the concatenated path with base path and have them created in advance.
func GetPathWithCreationWithBase(base, subdirs string) (string, error) {
	// Try to create the base directory first if not exists.
	if _, err := os.Stat(base); os.IsNotExist(err) {
		// Only create the top-most base directory.
		err := os.Mkdir(base, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Recursively create sub-directories.
	newPath := filepath.Join(base, subdirs)
	return newPath, os.MkdirAll(newPath, os.ModePerm)
}

// ReadEntireFile reads the entire file to a string.
func ReadEntireFile(fullPath string) string {
	b, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// GetFilePathListWithPattern returns the full paths for files matching the pattern in the base path.
func GetFilePathListWithPattern(basePath, fileNamePattern string) []string {
	var files []string
	filepath.WalkDir(basePath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if matched, _ := filepath.Match(fileNamePattern, d.Name()); matched {
			files = append(files, s)
		}
		return nil
	})
	log.Println("Found", len(files), "matched files with pattern:", fileNamePattern)
	return files
}

// LoadCookiesFile loads the external cookies file.
func LoadCookiesFile(filePath string) ([]*http.Cookie, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	return cookiestxt.Parse(f)
}

// LoadCookiesFileToString loads the external cookies file.
func LoadCookiesFileToString(filePath string) (string, error) {
	cookies, err := LoadCookiesFile(filePath)
	if err != nil {
		return "", nil
	}

	log.Println("Loaded", len(cookies), "cookies")
	c := make([]string, 0)
	for _, cookie := range cookies {
		if len(cookie.Name) == 0 || len(cookie.Value) == 0 {
			continue
		}
		c = append(c, fmt.Sprintf("%s=%s", cookie.Name, html.UnescapeString(cookie.Value)))
	}

	result := strings.Join(c, "; ")
	log.Println("Cookies are", result)
	return result, nil
}
