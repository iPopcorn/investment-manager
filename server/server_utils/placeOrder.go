package server_utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/util"
)

type PlaceOrderArgs struct {
	Client  *infrastructure.InvestmentManagerExternalHttpClient
	Offer   *types.Offer
	Preview bool
}

type previewOffer struct {
	ProductId         string                   `json:"product_id"`
	Side              types.Side               `json:"side"`
	Config            types.OrderConfiguration `json:"order_configuration"`
	RetailPortfolioId string                   `json:"retail_portfolio_id"`
}

func PlaceOrder(args *PlaceOrderArgs) ([]byte, error) {
	url := "https://api.coinbase.com/api/v3/brokerage/orders"
	var serializedRequest []byte
	var err error

	if args.Preview {
		log.Printf("Preview mode is on!\n")

		url += "/preview"

		previewReq := previewOffer{
			ProductId:         args.Offer.ProductId,
			Side:              args.Offer.Side,
			Config:            args.Offer.Config,
			RetailPortfolioId: args.Offer.RetailPortfolioId,
		}

		serializedRequest, err = json.Marshal(previewReq)

		if err != nil {
			log.Printf("Failed to serialize coinbase request\nrequest: %+v\nerror: %v\n", previewReq, err)
			return nil, err
		}
	} else {
		serializedRequest, err = json.Marshal(args.Offer)

		if err != nil {
			log.Printf("Failed to serialize coinbase request\nrequest: %+v\nerror: %v\n", args.Offer, err)
			return nil, err
		}
	}

	log.Printf("Sending create order request to coinbase\nreq: %s\n", string(serializedRequest))
	resp, err := args.Client.Post(url, serializedRequest)

	if err != nil {
		log.Printf("Error when sending req to coinbase\n%v\n", err)
		return nil, err
	}

	err = util.HandleErrorResponse(resp)

	if err != nil {
		fmt.Printf("Received error from coinbase while placing order\n%v\n", err)

		return nil, err
	}

	log.Printf("Received resp from coinbase\nresp: %s\n", string(resp))

	if args.Preview {
		var coinbaseResp types.CoinbaseOrderPreviewResponse

		err = json.Unmarshal(resp, &coinbaseResp)

		if err != nil {
			fmt.Printf("Failed to parse coinbase preview response\n%v\n", err)

			return nil, err
		}

		if len(coinbaseResp.Errors) > 0 {
			return nil, fmt.Errorf("Received errors from coinbase: %v", coinbaseResp.Errors)
		}
	} else {
		var coinbaseResp types.CoinbaseOrderPlacedResponse

		err = json.Unmarshal(resp, &coinbaseResp)

		if err != nil {
			fmt.Printf("Failed to parse coinbase response\n%v\n", err)

			return nil, err
		}
		return resp, nil
	}
	return resp, nil
}
