package abstractapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type AbstractApiClient struct {
	http   *http.Client
	config *Config
}

func NewAbstractApiClient(logger LoggerI, config *Config) *AbstractApiClient {
	if config == nil || len(config.ApiHosts) == 0 {
		logger.Fatalf("No api hosts specified")
	}
	config.apiHostsLatencies = make([]time.Duration, len(config.ApiHosts))

	if config.SortingType == RANDOM {
		rand.Seed(time.Now().UnixNano())
	}

	return &AbstractApiClient{
		http:   &http.Client{},
		config: config,
	}
}

func (c *AbstractApiClient) DoJsonReq(logger LoggerI, httpMethod string, method string, urlParams map[string]string, additionalHeaders map[string]string, output interface{}) error {
	var currIndex = -1
	var lastErr error
	for index := 0; index < len(c.config.ApiHosts); index++ {
		var host string
		if lastErr != nil {
			logger.Errorf("AbstractApiClient.DoJsonReq error: %v", lastErr)
		}
		if index > 0 {
			c.config.penaltyTheHost(currIndex)
			if httpMethod == "POST" {
				return lastErr
			}
		}
		var err error
		host, currIndex, err = c.config.getNextHost(currIndex)
		if err != nil {
			logger.Fatalf(err.Error())
		}

		secretPhrase, ok := urlParams["secretPhrase"]
		if ok {
			delete(urlParams, "secretPhrase")
		}
		logger.Debugf("Will request %v %v%v with params: %v", httpMethod, host, method, urlParams)
		if ok {
			urlParams["secretPhrase"] = secretPhrase
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
			lastErr = fmt.Errorf("couldn't unmarshal body of %v: %v. Body: %v", host+method, err, string(body))
			continue
		}
		latency := time.Since(startTime)
		c.config.appendLatencyToHost(latency, currIndex)

		return nil
	}
	c.config.penaltyTheHost(currIndex)
	return fmt.Errorf("couldn't get %v method: %v", method, lastErr)
}
