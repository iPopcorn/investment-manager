package types

type ExecuteStrategyRequest struct {
	Portfolio string `json:"portfolio"`
	Strategy  string `json:"strategy"`
	Currency  string `json:"currency"`
}

type Strategy struct {
	Name         string
	Currency     string
	OpenOffers   []Offer
	ClosedOffers []Offer
}

type Offer struct {
	ClientOrderId         string
	ProductId             string
	Side                  Side
	Config                OrderConfiguration
	SelfTradePreventionId string
	RetailPortfolioId     string
}

type OrderConfiguration struct {
	Type       OrderType
	BaseSize   string
	LimitPrice string
	PostOnly   bool
	EndTime    string // RFC3339 Timestamp
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
