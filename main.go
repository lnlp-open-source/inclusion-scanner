package main

import (
	"fmt"
	"github.com/lnlp-open-source/inclusion-scanner/lib/configuration"
	"github.com/lnlp-open-source/inclusion-scanner/lib/elasticsearch"
	"github.com/lnlp-open-source/inclusion-scanner/lib/filesystem"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

type InclusionCommand struct {
	cmd            *cobra.Command
	ConfigFilePath string
}

func main() {
	command := NewInclusionCommand()
	err := command.cmd.Execute()
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func NewInclusionCommand() *InclusionCommand {
	handler := InclusionCommand{}
	var cmd = &cobra.Command{
		Use: "inclusion-scanner",
		Example: `
Scan an individual project:
	./inclusion-scanner myproject/ 
Scan a parent directory with one or more project inside:
	./inclusion-scanner ~/myrepos`,
		Short: "Collects metrics on the usage of terms with racist connotations inside project directories",
		Long:  "Inclusion Scanner provides users with the ability to scan for racist and potentially non-inclusive terms and receive a report of findings so that they can be corrected with alternative terms.",
		Run:   handler.Run,
		Args:  cobra.MinimumNArgs(0),
	}
	handler.cmd = cmd
	cmd.Flags().StringVarP(&handler.ConfigFilePath, "config", "c", "config.yml", "File path to a config.yml file")
	return &handler
}

func (cmd *InclusionCommand) Run(_ *cobra.Command, args []string) {
	fmt.Println("Starting inclusion-scanner metrics collector... ")
	config, err := cmd.GetConfiguration()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = elasticsearch.DoHealthCheck(config.ElasticSearch.Url)
	if err != nil {
		fmt.Printf("ElasticSearch is not available at the URL '%s' specified in config.yml, %v.\n\nYou may need to run 'docker-compose up' and wait for the service to come online.\n", config.ElasticSearch.Url, err)
		return
	}

	err = startScans(config, args)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func startScans(config *configuration.Configuration, args []string) error {
	for _, directory := range args {
		fmt.Printf("Scanning file tree at %s\n", directory)
		fileSystemScanner := filesystem.NewFileSystemScanner(config)
		err := fileSystemScanner.ScanDirectory(directory)
		if err != nil {
			return fmt.Errorf("Could not collect racist terms. %v", err)
		}
	}
	return nil
}

func (cmd *InclusionCommand) GetConfiguration() (*configuration.Configuration, error) {
	if _, err := os.Stat(cmd.ConfigFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Config file at path '%s' does not exist\n", cmd.ConfigFilePath)
	}
	configFileBytes, err := ioutil.ReadFile(cmd.ConfigFilePath)
	if err != nil {
		return nil, err
	}
	config, err := configuration.NewConfigurationFromBytes(configFileBytes)
	if err != nil {
		return nil, err
	}
	err = config.CheckValidity()
	if err != nil {
		return nil, err
	}
	return config, nil
}
