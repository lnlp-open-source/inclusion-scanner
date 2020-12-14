package elasticsearch

import (
	"fmt"
	"net/http"
)

func DoHealthCheck(url string) error {
	healthEndpoint := "/_cluster/health"
	if url[0] == '/' {
		healthEndpoint = healthEndpoint[1:]
	}
	url = url + healthEndpoint
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("health endpoint returned status code %d", response.StatusCode)
	}
	return nil
}
