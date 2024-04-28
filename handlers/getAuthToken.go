package handlers

import (
	"fmt"

	"github.com/iPopcorn/investment-manager/auth"
	"github.com/spf13/cobra"
)

func HandleGetAuthToken(cmd *cobra.Command, args []string) error {
	apiKey, err := auth.GetApiKey()
	httpMethod := "GET"
	host := "api.coinbase.com"
	path := "/api/v3/brokerage/products"

	jwtOptions := auth.BuildJWTOptions{
		Service:    "retail_rest_api_proxy",
		Uri:        fmt.Sprintf("%s %s%s", httpMethod, host, path),
		PrivateKey: apiKey.PrivateKey,
		Name:       apiKey.Name,
	}

	token, err := auth.BuildJWT(jwtOptions)

	if err != nil {
		fmt.Printf("error getting jwt\n%v\n", err)
		return err
	}

	fmt.Printf("Token:\n%s\n", token)
	return nil
}
