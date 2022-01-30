package abstractapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type AbstractApiClient struct {
	http          *http.Client
	ApiHost       string
	staticHeaders map[string]string
}

func NewAbstractApiClient(apiHost string, staticHeaders map[string]string) *AbstractApiClient {
	return &AbstractApiClient{
		http:          &http.Client{},
		ApiHost:       apiHost,
		staticHeaders: staticHeaders,
	}
}

func (c *AbstractApiClient) DoJsonReq(logger LoggerI, httpMethod string, method string, urlParams map[string]string, additionalHeaders map[string]string, output interface{}) error {
	// protect sensitive data from logging
	var sensitiveData = map[string]string{}
	for _, key := range []string{"secretPhrase", "messageToEncrypt"} {
		data, ok := urlParams[key]
		if ok {
			sensitiveData[key] = data
			delete(urlParams, key)
		}
	}
	logger.Debugf("Request %v %v%v with params: %v", httpMethod, c.ApiHost, method, urlParams)
	for key, data := range sensitiveData {
		urlParams[key] = data
	}

	// requesting
	req, err := http.NewRequest(httpMethod, c.ApiHost+method, nil)
	if err != nil {
		return fmt.Errorf("error create req %v", c.ApiHost+method)
	}

	req.Header.Set("Accepts", "application/json")
	if c.staticHeaders != nil {
		for key, value := range c.staticHeaders {
			req.Header.Add(key, value)
		}
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
		return fmt.Errorf("error perform %v %v: %v", httpMethod, c.ApiHost+method, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't read body of %v: %v", c.ApiHost+method, err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error StatusCode %v for %v. Body: %v", resp.StatusCode, c.ApiHost+method, strconv.Quote(string(body)))
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal body of %v: %v. Body: %v", c.ApiHost+method, err, strconv.Quote(string(body)))
	}

	return nil
}
