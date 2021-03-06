package filesystem

import (
	"fmt"
	"github.com/lnlp-open-source/inclusion-scanner/lib/configuration"
	"github.com/lnlp-open-source/inclusion-scanner/lib/elasticsearch"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileSystemScanner struct {
	Config *configuration.Configuration
}

func NewFileSystemScanner(config *configuration.Configuration) *FileSystemScanner {
	return &FileSystemScanner{Config: config}
}

func unique(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}
	for e := range arr {
		if !occured[arr[e]] {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func (scanner *FileSystemScanner) ScanDirectory(directoryPath string) error {
	return filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if !scanner.ShouldScanDirectory(path) {
				return filepath.SkipDir
			}
			return nil
		}
		if !scanner.ShouldScanFileAtPath(path) {
			return nil
		}

		nonInclusiveTermsUsed, scanError := scanner.ScanFileAtPath(path)
		if scanError != nil {
			return fmt.Errorf("Failed to scan file, %v", scanError)
		}
		if len(nonInclusiveTermsUsed) > 0 {
			fmt.Printf("File %s contains non-inclusive terms: %s\n", path, strings.Join(unique(nonInclusiveTermsUsed), ", "))
			return elasticsearch.StoreScan(scanner.Config, path, nonInclusiveTermsUsed)
		}
		return nil
	})
}

func (scanner *FileSystemScanner) ScanFileAtPath(filePath string) (nonInclusiveTermsUsed []string, err error) {
	// NOTE: This allows duplicates, such as "master master" if the term is used more than once
	if !scanner.ShouldScanFileAtPath(filePath) {
		return
	}
	if _, statError := os.Stat(filePath); os.IsNotExist(statError) {
		return
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	fileContent := strings.ToLower(string(fileBytes))
	re := regexp.MustCompile("(?i)((" + strings.Join(scanner.Config.Terms, ")|(") + "))")
	nonInclusiveTermsUsed = append(nonInclusiveTermsUsed, re.FindAllString(fileContent, -1)...)
	return
}

func (scanner *FileSystemScanner) ShouldScanFileAtPath(filePath string) bool {
	valid := false
	for _, extension := range scanner.Config.IncludedExtensions {
		if strings.Contains(filePath, extension) {
			valid = true
		}
	}
	return valid
}

func (scanner *FileSystemScanner) ShouldScanDirectory(directoryPath string) bool {
	valid := true
	for _, excludedDirectory := range scanner.Config.ExcludedDirectories {
		if strings.Contains(directoryPath, excludedDirectory) {
			valid = false
		}
	}
	return valid
}
