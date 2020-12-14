package configuration

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type Configuration struct {
	Scanners      ScannerList         `yaml:"scanners"`
	ElasticSearch ElasticSearchConfig `yaml:"elasticsearch"`

	Terms               []string `yaml:"terms"`
	IncludedExtensions  []string `yaml:"included_extensions"`
	ExcludedDirectories []string `yaml:"excluded_directories"`
}

type ScannerList struct {
	Repositories RepositoryConfig `yaml:"repositories"`
}

type RepositoryConfig struct {
	Index string `yaml:"database_index"`
}

type ElasticSearchConfig struct {
	Url string `yaml:"url"`
}

func NewConfigurationFromBytes(configBytes []byte) (*Configuration, error) {
	config := &Configuration{}
	err := yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	config.ElasticSearch.Url = formatUrl(config.ElasticSearch.Url)
	//_, err = url.ParseRequestURI(config.ElasticSearch.Url)
	//if err != nil {
	//	return nil, fmt.Errorf("ElasticSearch URL '%s' in config.yml is not in a URL format.", config.ElasticSearch.Url)
	//}
	return config, nil
}

func formatUrl(url string) string {
	url = addHttpPrefix(url)
	return removeTrailingSlash(url)
}

func addHttpPrefix(url string) string {
	if len(url) == 0 {
		return url
	}

	httpPrefix := "http://"
	httpsPrefix := "https://"
	if !strings.HasPrefix(url, httpPrefix) && !strings.HasPrefix(url, httpsPrefix) {
		url = httpPrefix + url
	}
	return url
}

func removeTrailingSlash(url string) string {
	if url != "" && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	return url
}

func (config *Configuration) CheckValidity() error {
	if len(config.ElasticSearch.Url) == 0 {
		return fmt.Errorf("An elasticsearch url is required in the config.yml file.")
	} else if config.Terms == nil || len(config.Terms) == 0 {
		return fmt.Errorf("One or more terms are required in the config.yml file.")
	} else if config.IncludedExtensions == nil || len(config.IncludedExtensions) == 0 {
		return fmt.Errorf("One or more included file extensions are required in the config.yml file.")
	}
	return nil
}
