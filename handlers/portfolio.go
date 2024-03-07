package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/iPopcorn/investment-manager/auth"
	"github.com/spf13/cobra"
)

type PortfolioResponse struct {
	Portfolios []Portfolio `json:"portfolios"`
}

type Portfolio struct {
	Name    string `json:"name"`
	Uuid    string `json:"uuid"`
	Type    string `json:"type"`
	Deleted bool   `json:"deleted"`
}

func HandlePortfolio(cmd *cobra.Command, args []string) {
	fmt.Printf("portfolio called\nargs: %v\n", args)
	fmt.Printf("listing portfolios...\n")
	listPortfolios()
}

func listPortfolios() {
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"
	host := "api.coinbase.com"
	path := "/api/v3/brokerage/portfolios"
	method := "GET"

	apiKey, err := auth.GetApiKey()
	if err != nil {
		fmt.Printf("error getting API Key\n%v\n", err)
	}

	jwtOptions := auth.BuildJWTOptions{
		Service:    "retail_rest_api_proxy",
		Uri:        fmt.Sprintf("%s %s%s", method, host, path),
		PrivateKey: apiKey.PrivateKey,
		Name:       apiKey.Name,
	}
	token, err := auth.BuildJWT(jwtOptions)

	if err != nil {
		fmt.Printf("error getting jwt\n%v\n", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	authHeader := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var resp PortfolioResponse
	json.Unmarshal(body, &resp)

	fmt.Println("Portfolios:")
	for _, p := range resp.Portfolios {
		fmt.Printf("Name: %s\n", p.Name)
		fmt.Printf("Type: %s\n", p.Type)
		fmt.Printf("UUID: %s\n", p.Uuid)
		fmt.Printf("Is Deleted?: %t\n", p.Deleted)
	}
}
