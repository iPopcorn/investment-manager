package types

type ProductResponse struct {
	Products []Product `json:"products"`
}

// TODO: Add fields as needed (https://docs.cloud.coinbase.com/advanced-trade/reference/retailbrokerageapi_getproducts)
type Product struct {
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
}
