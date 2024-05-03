package util

import (
	"encoding/json"
	"fmt"

	"github.com/iPopcorn/investment-manager/types"
)

func HandleErrorResponse(resp []byte) error {
	var errResp types.ErrorResponse

	err := json.Unmarshal(resp, &errResp)

	if err != nil {
		fmt.Printf("Failed to parse response.\n%q\n", string(resp))
		return err
	}

	if errResp.Error != "" {
		return fmt.Errorf("Error: %s\nMessage: %s\n", errResp.Error, errResp.Message)
	}

	return nil
}
