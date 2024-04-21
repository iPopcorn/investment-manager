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
