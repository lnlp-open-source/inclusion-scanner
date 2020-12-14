package configuration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigurationFromBytes(t *testing.T) {
	t.Run("ElasticSearch URL loaded and trailing slash automatically removed", func(t *testing.T) {
		config, err := NewConfigurationFromBytes([]byte("elasticsearch:\n  url: http://localhost:9200/"))
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost:9200", config.ElasticSearch.Url)
	})

	t.Run("Exact number of terms loaded", func(t *testing.T) {
		config, err := NewConfigurationFromBytes([]byte("elasticsearch:\n  url: http://localhost:9200/\nterms:\n  - \"master\"\n  - \"slave\""))
		assert.NoError(t, err)
		assert.Equal(t, 2, len(config.Terms))
	})

	t.Run("Exact number of included extensions loaded", func(t *testing.T) {
		config, err := NewConfigurationFromBytes([]byte("elasticsearch:\n  url: http://localhost:9200/\nincluded_extensions:\n  - \".yaml\""))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(config.IncludedExtensions))
	})

	t.Run("Exact number of excluded directories loaded", func(t *testing.T) {
		config, err := NewConfigurationFromBytes([]byte("elasticsearch:\n  url: http://localhost:9200/\nexcluded_directories:\n  - \".git\"\n  - \".DS_Store\""))
		assert.NoError(t, err)
		assert.Equal(t, 2, len(config.ExcludedDirectories))
	})
}

func TestConfigurationIsValid(t *testing.T) {
	t.Run("Config is missing an elasticsearch url", func(t *testing.T) {
		config := &Configuration{}
		assert.EqualError(t, config.CheckValidity(), "An elasticsearch url is required in the config.yml file.")
	})

	t.Run("Config does not have any terms", func(t *testing.T) {
		config := &Configuration{ElasticSearch: ElasticSearchConfig{Url: "http://localhost:9200"}}
		assert.EqualError(t, config.CheckValidity(), "One or more terms are required in the config.yml file.")
	})

	t.Run("Config does not have any included extensions", func(t *testing.T) {
		config := &Configuration{ElasticSearch: ElasticSearchConfig{Url: "http://localhost:9200"}, Terms: []string{"blacklist"}}
		assert.EqualError(t, config.CheckValidity(), "One or more included file extensions are required in the config.yml file.")
	})

	t.Run("Config is valid", func(t *testing.T) {
		config := &Configuration{ElasticSearch: ElasticSearchConfig{Url: "http://localhost:9200"}, Terms: []string{"blacklist"}, IncludedExtensions: []string{".go"}}
		assert.NoError(t, config.CheckValidity())
	})
}
