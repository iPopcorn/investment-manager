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

	senderName := args[0]
	receiverName := args[1]

	portfolios, err := listPortfolios(internalClient)
	if err != nil {
		return fmt.Errorf("could not get portfolios\n%v\n", err)
	}

	senderID, err := getPortfolioIdByName(senderName, portfolios)

	if err != nil {
		return fmt.Errorf("Could not get sender ID\n%v\n", err)
	}

	receiverID, err := getPortfolioIdByName(receiverName, portfolios)

	if err != nil {
		return fmt.Errorf("Could not receiver ID\n%v\n", err)
	}

	request := types.TransferRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     args[2],
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

func getPortfolioIdByName(name string, portfolios *types.PortfolioResponse) (string, error) {
	for _, portfolio := range portfolios.Portfolios {
		if name == portfolio.Name {
			return portfolio.Uuid, nil
		}
	}

	return "", fmt.Errorf("Could not find portfolio id given %q\nportfolios: %v\n", name, portfolios)
}
