package types

type ProductResponse struct {
	Products []Product `json:"products"`
}

// TODO: Add fields as needed (https://docs.cloud.coinbase.com/advanced-trade/reference/retailbrokerageapi_getproducts)
type Product struct {
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
}

type BestBidAskResponse struct {
	PriceBooks []PriceBook `json:"pricebooks"`
}

type PriceBook struct {
	ProductID string `json:"product_id"`
	Bids      []Bid  `json:"bids"`
	Asks      []Bid  `json:"asks"`
	Time      string `json:"time"`
}

type Bid struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}
