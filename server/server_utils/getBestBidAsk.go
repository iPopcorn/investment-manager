package server_utils

import (
	"fmt"
	"log"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/types/mappers"
)

func GetBestBidAsk(
	client *infrastructure.InvestmentManagerExternalHttpClient,
	productID string,
) (*types.BestBidAskResponse, error) {
	url := fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/best_bid_ask?product_ids=%s", productID)
	resp, err := client.Get(url)

	log.Printf("raw resp: \n%s\n", string(resp))

	if err != nil {
		log.Printf("Error retrieving best bid ask from URL: %q\nError: %v\n", url, err)
		return nil, err
	}

	bestBidAskResponse, err := mappers.MapBestBidAskResponse(resp)

	if err != nil {
		return nil, fmt.Errorf("Failed to map product response to object\n%v\n", err)
	}

	if len(bestBidAskResponse.PriceBooks) != 1 {
		return nil, fmt.Errorf("Invalid response, price books has unexpected length: %d\n", len(bestBidAskResponse.PriceBooks))
	}

	if len(bestBidAskResponse.PriceBooks[0].Bids) < 1 {
		return nil, fmt.Errorf("Invalid price book, no bids in list\nProductID: %q\n", productID)
	}

	if len(bestBidAskResponse.PriceBooks[0].Asks) < 1 {
		return nil, fmt.Errorf("Invalid price book, no asks in list\nProductID: %q\n", productID)
	}

	return bestBidAskResponse, nil
}
