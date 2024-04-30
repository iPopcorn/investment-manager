package types

type ExecuteStrategyRequest struct {
	Portfolio string `json:"portfolio"`
	Strategy  string `json:"strategy"`
	Currency  string `json:"currency"`
}

type Strategy struct {
	Name         string  `json:"name"`
	Currency     string  `json:"currency"`
	OpenOffers   []Offer `json:"open_offers"`
	ClosedOffers []Offer `json:"closed_offers"`
}

type Offer struct {
	ClientOrderId         string             `json:"client_order_id"`
	ProductId             string             `json:"product_id"`
	Side                  Side               `json:"side"`
	Config                OrderConfiguration `json:"config"`
	SelfTradePreventionId string             `json:"self_trade_prevention_id"`
	RetailPortfolioId     string             `json:"retail_portfolio_id"`
}

type OrderConfiguration struct {
	Type       OrderType `json:"type"`
	BaseSize   string    `json:"base_size"`   // Amount of base currency to spend on order
	LimitPrice string    `json:"limit_price"` // Ceiling price for which the order should get filled.
	PostOnly   bool      `json:"post_only"`
	EndTime    string    `json:"end_time"` // RFC3339 Timestamp
}

type OrderType string

const (
	LimitLimitGTD OrderType = "limit_limit_gtd"
)

type Side string

const (
	BUY  Side = "BUY"
	SELL Side = "SELL"
)
