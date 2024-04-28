package mappers

import (
	"encoding/json"
	"fmt"

	"github.com/iPopcorn/investment-manager/types"
)

func MapPortfolioResponse(httpResponse []byte) (*types.PortfolioResponse, error) {
	var resp types.PortfolioResponse
	err := json.Unmarshal(httpResponse, &resp)

	if err != nil {
		return nil, fmt.Errorf("Failed to map http response to object\n%v", err)
	}

	return &resp, nil
}

func MapPortfolioDetailsResponse(httpResponse []byte) (*types.PortfolioDetailsResponse, error) {
	var resp types.PortfolioDetailsResponse
	err := json.Unmarshal(httpResponse, &resp)

	if err != nil {
		return nil, fmt.Errorf("Failed to map http response to object\n%v", err)
	}

	if resp.Breakdown.Portfolio.Name == "" {
		return nil, fmt.Errorf("Failed to map Portfolio details from response\nGiven: %q\n", string(httpResponse))
	}

	return &resp, nil
}
