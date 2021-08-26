package abstractapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type AbstractApiClient struct {
	http   *http.Client
	config *Config
}

func NewAbstractApiClient(config *Config) *AbstractApiClient {
	if config == nil || len(config.ApiHosts) == 0 {
		log.Fatalf("No api hosts specified")
	}
	config.apiHostsLatencies = make([]time.Duration, len(config.ApiHosts))

	return &AbstractApiClient{
		http:   &http.Client{},
		config: config,
	}
}

func (c *AbstractApiClient) DoJsonReq(httpMethod string, method string, urlParams map[string]string, additionalHeaders map[string]string, output interface{}) error {
	var bestIndex int
	var lastErr error
	for index := 0; index < len(c.config.ApiHosts); index++ {
		var host string
		if lastErr != nil && c.config.Debug {
			log.Printf("AbstractApiClient.DoJsonReq error: %v", lastErr)
		}
		if index > 0 {
			c.config.penaltyTheHost(bestIndex)
			if httpMethod == "POST" {
				return lastErr
			}
		}
		host, bestIndex = c.config.getBestHost()

		if c.config.Debug {
			secretPhrase, ok := urlParams["secretPhrase"]
			if ok {
				delete(urlParams, "secretPhrase")
			}
			log.Printf("Will request %v %v%v with params: %v", httpMethod, host, method, urlParams)
			if ok {
				urlParams["secretPhrase"] = secretPhrase
			}
		}

		req, err := http.NewRequest(httpMethod, host+method, nil)
		if err != nil {
			lastErr = fmt.Errorf("error create req %v", host+method)
			continue
		}

		req.Header.Set("Accepts", "application/json")
		for key, value := range c.config.StaticHeaders {
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

		startTime := time.Now()
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
		latency := time.Since(startTime)
		c.config.appendLatencyToHost(latency, bestIndex)

		return nil
	}
	c.config.penaltyTheHost(bestIndex)
	return fmt.Errorf("couldn't get %v method: %v", method, lastErr)
}
