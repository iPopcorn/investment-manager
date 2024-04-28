package mappers

import (
	"encoding/json"
	"fmt"

	"github.com/iPopcorn/investment-manager/types"
)

func MapProductResponse(httpResponse []byte) (*types.ProductResponse, error) {
	var resp types.ProductResponse
	err := json.Unmarshal(httpResponse, &resp)

	if err != nil {
		return nil, fmt.Errorf("Failed to map http response to object\n%v", err)
	}

	return &resp, nil
}
