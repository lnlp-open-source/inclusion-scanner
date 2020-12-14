package filesystem

import (
	"github.com/lnlp-open-source/inclusion-scanner/lib/configuration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldScanFileAtPath(t *testing.T) {
	t.Run("File has extension that is inside included_extensions configuration", func(t *testing.T) {
		scanner := NewFileSystemScanner(&configuration.Configuration{
			IncludedExtensions: []string{".go"},
		})
		assert.True(t, scanner.ShouldScanFileAtPath("myfile.go"))
	})

	t.Run("File has a name that is inside included_extensions configuration", func(t *testing.T) {
		scanner := NewFileSystemScanner(&configuration.Configuration{
			IncludedExtensions: []string{"Jenkinsfile"},
		})
		assert.True(t, scanner.ShouldScanFileAtPath("Jenkinsfile"))
	})

	t.Run("File has extension that is not in included_extensions configuration", func(t *testing.T) {
		scanner := NewFileSystemScanner(&configuration.Configuration{
			IncludedExtensions: []string{".go"},
		})
		assert.False(t, scanner.ShouldScanFileAtPath("myfile.sh"))
	})
}

func TestShouldScanDirectory(t *testing.T) {
	t.Run("Directory has extension that is in excluded_directories configuration", func(t *testing.T) {
		scanner := NewFileSystemScanner(&configuration.Configuration{
			ExcludedDirectories: []string{"somedir"},
		})
		assert.False(t, scanner.ShouldScanDirectory("myproject/somedir"))
	})

	t.Run("Directory has extension that is not in excluded_directories configuration", func(t *testing.T) {
		scanner := NewFileSystemScanner(&configuration.Configuration{
			ExcludedDirectories: []string{"somedir"},
		})
		assert.True(t, scanner.ShouldScanDirectory("myproject/mydir"))
	})
}
