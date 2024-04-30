package infrastructure

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/iPopcorn/investment-manager/auth"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type InvestmentManagerExternalHttpClient struct {
	HttpClient HttpClient
}

func GetInvestmentManagerExternalHttpClient() *InvestmentManagerExternalHttpClient {
	return &InvestmentManagerExternalHttpClient{
		HttpClient: &http.Client{},
	}
}

func (client InvestmentManagerExternalHttpClient) Get(url string) ([]byte, error) {
	return client.sendAuthenticatedHttpRequest(url, "GET", nil)
}

func (client InvestmentManagerExternalHttpClient) Post(url string, request []byte) ([]byte, error) {
	return client.sendAuthenticatedHttpRequest(url, "POST", request)
}

func (client InvestmentManagerExternalHttpClient) sendAuthenticatedHttpRequest(url, method string, request []byte) ([]byte, error) {
	emptyResponse := []byte{}

	jwt, err := getJWT(url, method)
	if err != nil {
		fmt.Println("Failed to get jwt")
		return emptyResponse, err
	}

	var req *http.Request
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

	authHeader := fmt.Sprintf("Bearer %s", jwt)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	res, err := client.HttpClient.Do(req)

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

func getJWT(url, httpMethod string) (string, error) {
	apiKey, err := auth.GetApiKey()
	if err != nil {
		fmt.Printf("error getting API Key\n%v\n", err)
		return "", err
	}

	// separate "https://" from the rest of the url
	token1 := strings.Split(url, "//")

	if len(token1) != 2 {
		return "", errors.New(fmt.Sprintf("Invalid url: %s", url))
	}

	// separate domain from rest of the path token1[0] is expected to equal "https:"
	token2 := strings.Split(token1[1], "/")

	// domain is expected to be similar to "api.coinbase.com"
	token3 := strings.Split(token2[0], ".")

	if len(token3) < 2 && len(token3) > 3 {
		return "", errors.New(fmt.Sprintf("Invalid url: %s", url))
	}

	host := token2[0]
	var path string

	for i, token := range token2 {
		// skip first index of token2 because it is the host
		if i > 0 {
			path += "/" + token
		}
	}

	// remove query params
	path = strings.Split(path, "?")[0]
	uri := fmt.Sprintf("%s %s%s", httpMethod, host, path)

	jwtOptions := auth.BuildJWTOptions{
		Service:    "retail_rest_api_proxy",
		Uri:        uri,
		PrivateKey: apiKey.PrivateKey,
		Name:       apiKey.Name,
	}

	token, err := auth.BuildJWT(jwtOptions)

	if err != nil {
		fmt.Printf("error getting jwt\n%v\n", err)
		return "", err
	}

	return token, nil
}
