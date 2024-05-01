package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/spf13/cobra"
)

func TransferFundsHandlerFactory(internalClient *infrastructure.InvestmentManagerInternalHttpClient) CobraCommandHandler {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("Unexpected number of args.\nExpected 3, Received %d", len(args))
		}
		err := transferFundsHandler(internalClient, args)

		return err
	}
}

func transferFundsHandler(internalClient *infrastructure.InvestmentManagerInternalHttpClient, args []string) error {
	_, err := strconv.ParseFloat(args[2], 64)

	if err != nil {
		return fmt.Errorf("Could not convert arg to float.\nGiven: %s\n%v\n", args[2], err)
	}

	request := types.TransferRequest{
		Sender:   args[0],
		Receiver: args[1],
		Amount:   args[2],
	}

	serializedRequest, err := json.Marshal(request)

	if err != nil {
		return fmt.Errorf("Failed to serialize request\n%v\n", err)
	}

	resp, err := internalClient.Post("/transfer-funds", serializedRequest)

	if err != nil {
		return fmt.Errorf("Response failed: %v\n", err)
	}

	fmt.Printf("Success!\n%s\n", string(resp))
	return nil
}
