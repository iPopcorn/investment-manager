package util

import (
	"fmt"
	"log"
	"strings"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/types/mappers"
)

func GetProductID(
	client *infrastructure.InvestmentManagerExternalHttpClient,
	portfolioDetails *types.PortfolioDetailsResponse,
	baseCurrency string,
) (string, error) {
	quoteCurrencyID := portfolioDetails.Breakdown.PortfolioBalances.TotalCashEquivalentBalance.Currency

	productID := strings.ToUpper(baseCurrency) + "-" + quoteCurrencyID

	url := fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/products?product_type=SPOT&product_ids=%s", productID)
	resp, err := client.Get(url)

	if err != nil {
		log.Printf("Error retrieving product list from URL: %q\nError: %v", url, err)
		return "", err
	}

	productResponse, err := mappers.MapProductResponse(resp)

	if err != nil {
		return "", fmt.Errorf("Failed to map product response to object\n%v\n", err)
	}

	if len(productResponse.Products) < 1 {
		return "", fmt.Errorf("Invalid product id, no products in list\nProductID: %q", productID)
	}

	if len(productResponse.Products) > 1 {
		fmt.Println("Warning, found multiple products for given id")
	}

	return productID, nil
}
