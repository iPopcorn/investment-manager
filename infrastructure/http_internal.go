package infrastructure

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type InvestmentManagerInternalHttpClient struct {
	client  HttpClient
	baseURL string
}

func GetInvestmentManagerInternalHttpClient() *InvestmentManagerInternalHttpClient {
	return &InvestmentManagerInternalHttpClient{
		client:  &http.Client{},
		baseURL: "http://127.0.0.1:5000",
	}
}

func (c *InvestmentManagerInternalHttpClient) Get(path string) ([]byte, error) {
	url := c.baseURL + path
	return c.sendInternalHttpRequest(url, "GET", nil)
}

func (c *InvestmentManagerInternalHttpClient) sendInternalHttpRequest(url, method string, request []byte) ([]byte, error) {
	emptyResponse := []byte{}

	var req *http.Request
	var err error
	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else if method == "POST" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(request))
	} else {
		return emptyResponse, errors.New(fmt.Sprintf("Unsupported http verb, recieved %s", method))
	}

	if err != nil {
		fmt.Println("Failed to build request")
		return emptyResponse, err
	}

	res, err := c.client.Do(req)

	if err != nil {
		fmt.Println("Failed to get response")
		return emptyResponse, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Failed to read response")
		return emptyResponse, err
	}

	return body, nil
}
