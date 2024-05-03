package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/server/server_utils"
)

type HandlePortfolioArgs struct {
	Client *infrastructure.InvestmentManagerExternalHttpClient
	Writer http.ResponseWriter
	Req    *http.Request
	Args   []string
}

func HandlePortfolio(hpArgs HandlePortfolioArgs) {
	w := hpArgs.Writer
	client := hpArgs.Client
	r := hpArgs.Req

	w.Header().Set("Content-Type", "application/json")
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"

	if r.Method == http.MethodPost {
		body := r.Body

		defer body.Close()

		bodyData, err := ioutil.ReadAll(body)

		if err != nil {
			log.Printf("Failed to read body from request: %v\n", err)
			server_utils.WriteResponse(w, nil, err)
		}

		resp, err := client.Post(url, bodyData)

		server_utils.WriteResponse(w, resp, err)
	} else {
		if len(hpArgs.Args) == 1 {
			portfolioUUID := hpArgs.Args[0]
			url = url + "/" + portfolioUUID
			resp, err := client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolio details from URL: %q\nError: %v", url, err)
			}

			server_utils.WriteResponse(w, resp, err)
		} else {
			resp, err := client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolios from URL: %q\nError: %v", url, err)
			}

			server_utils.WriteResponse(w, resp, err)
		}
	}

	log.Printf("Request handled successfully!")
}
