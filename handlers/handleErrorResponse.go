package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/iPopcorn/investment-manager/types"
)

func handleErrorResponse(resp []byte) error {
	var errResp types.ErrorResponse

	err := json.Unmarshal(resp, &errResp)

	if err != nil {
		fmt.Println("Failed to parse response.")
		return err
	}

	if errResp.Error != "" {
		errMsg := fmt.Sprintf("Error: %s\nMessage: %s",
			errResp.Error,
			errResp.Message,
		)

		return errors.New(errMsg)
	}

	return nil
}
