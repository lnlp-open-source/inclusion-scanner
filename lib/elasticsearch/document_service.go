package elasticsearch

import (
	"bytes"
	"fmt"
	"github.com/lnlp-open-source/inclusion-scanner/lib/configuration"
	"io/ioutil"
	"net/http"
	"time"
)

func StoreScan(config *configuration.Configuration, filePath string, nonInclusiveTermsUsed []string) error {
	document := NewDocument(filePath, nonInclusiveTermsUsed)
	bodyPayload, err := document.GetPayload()
	if err != nil {
		return err
	}

	index := config.Scanners.Repositories.Index + "-" + time.Now().Format("2006.01.02")
	httpClient := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/_doc/%s", config.ElasticSearch.Url, index, document.Id), bytes.NewReader(bodyPayload))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	responseContent := ""
	if response.Body != nil {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		responseContent = string(responseBytes)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusCreated {
		if len(responseContent) > 0 {
			return fmt.Errorf("Could not store file results in elasticsearch, http status is %d.\n%s", response.StatusCode, responseContent)
		}
		return fmt.Errorf("Could not store file results in elasticsearch, http status is %d.", response.StatusCode)
	}
	return nil
}
