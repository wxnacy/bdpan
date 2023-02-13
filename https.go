package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
)

func Do2HTTPRequest(url string, body io.Reader, headers map[string]string) (string, int, error) {
	// timeout := 500 * time.Second
	retryTimes := 3
	tr := &http.Transport{
		MaxIdleConnsPerHost: -1,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	// httpClient.Timeout = timeout
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", 0, err
	}
	// request header
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var resp *http.Response
	for i := 1; i <= retryTimes; i++ {
		resp, err = httpClient.Do(req)
		if err == nil {
			break
		}
		if i == retryTimes {
			return "", 0, err
		}
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	return string(respBody), resp.StatusCode, nil
}

func Do2HTTPRequestToBytes(url string, body io.Reader, headers map[string]string) ([]byte, int, error) {
	// timeout := 500 * time.Second
	retryTimes := 3
	tr := &http.Transport{
		MaxIdleConnsPerHost: -1,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	// httpClient.Timeout = timeout
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, 0, err
	}
	// request header
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var resp *http.Response
	for i := 1; i <= retryTimes; i++ {
		resp, err = httpClient.Do(req)
		if err == nil {
			break
		}
		if i == retryTimes {
			return nil, 0, err
		}
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respBody, resp.StatusCode, nil
}
