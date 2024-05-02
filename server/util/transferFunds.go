package util

import (
	"encoding/json"
	"log"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
)

type coinbaseTransferFundsRequest struct {
	Funds      Funds  `json:"funds"`
	SenderID   string `json:"source_portfolio_uuid"`
	ReceiverID string `json:"target_portfolio_uuid"`
}

type Funds struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

func TransferFunds(client *infrastructure.InvestmentManagerExternalHttpClient, req *types.TransferRequest) ([]byte, error) {
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios/move_funds"

	coinbaseReq := coinbaseTransferFundsRequest{
		Funds: Funds{
			Value:    req.Amount,
			Currency: "GBP",
		},
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
	}

	serializedReq, err := json.Marshal(coinbaseReq)
	if err != nil {
		log.Printf("Failed to serialize request\nGiven: %+v\n", req)
		return nil, err
	}

	log.Printf("Sending request to coinbase: %q\n", string(serializedReq))
	resp, err := client.Post(url, serializedReq)

	return resp, err
}
