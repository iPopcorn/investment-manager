package infrastructure

import (
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

type InvestmentManagerHTTPClient struct {
	HttpClient HttpClient
}

func (client InvestmentManagerHTTPClient) Get(url string) ([]byte, error) {
	httpMethod := "GET"
	emptyResponse := []byte{}

	jwt, err := getJWT(url, httpMethod)
	if err != nil {
		return emptyResponse, err
	}

	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		fmt.Println(err)
		return emptyResponse, err
	}

	authHeader := fmt.Sprintf("Bearer %s", jwt)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	res, err := client.HttpClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return emptyResponse, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
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

	jwtOptions := auth.BuildJWTOptions{
		Service:    "retail_rest_api_proxy",
		Uri:        fmt.Sprintf("%s %s%s", httpMethod, host, path),
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
