package abstractapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type AbstractApiClient struct {
	http *http.Client

	// config
	Debug         bool
	apiHosts      []string
	staticHeaders map[string]string
}

func NewAbstractApiClient(apiHosts []string, staticHeaders map[string]string, debug bool) *AbstractApiClient {
	return &AbstractApiClient{
		http:          &http.Client{},
		Debug:         debug,
		apiHosts:      apiHosts,
		staticHeaders: staticHeaders,
	}
}

func (c *AbstractApiClient) DoJsonReq(httpMethod string, method string, urlParams map[string]string, additionalHeaders map[string]string, output interface{}) error {
	if c.Debug {
		secretPhrase, ok := urlParams["secretPhrase"]
		if ok {
			delete(urlParams, "secretPhrase")
		}
		log.Printf("Will request %v %v with params: %v", httpMethod, method, urlParams)
		if ok {
			urlParams["secretPhrase"] = secretPhrase
		}
	}

	var lastErr error
	for index, host := range c.apiHosts {
		if lastErr != nil && c.Debug {
			log.Printf("AbstractApiClient.makeJsonReq error: %v", lastErr)
		}
		if index > 0 && httpMethod == "POST" {
			return lastErr
		}

		req, err := http.NewRequest(httpMethod, host+method, nil)
		if err != nil {
			lastErr = fmt.Errorf("error create req %v", host+method)
			continue
		}

		req.Header.Set("Accepts", "application/json")
		for key, value := range c.staticHeaders {
			req.Header.Add(key, value)
		}
		for key, value := range additionalHeaders {
			req.Header.Add(key, value)
		}

		q := url.Values{}
		for key, value := range urlParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error perform %v %v: %v", httpMethod, host+method, err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("couldn't read body of %v: %v", host+method, err)
			continue
		}

		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("error StatusCode %v for %v: %v", resp.StatusCode, host+method, string(body))
			continue
		}

		err = json.Unmarshal(body, output)
		if err != nil {
			lastErr = fmt.Errorf("couldn't unmarshal body of %v: %v", host+method, err)
			continue
		}
		return nil
	}
	return fmt.Errorf("couldn't get %v method: %v", method, lastErr)
}
