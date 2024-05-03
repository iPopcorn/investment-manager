package types

type ExecuteStrategyRequest struct {
	Portfolio string            `json:"portfolio"`
	Strategy  StrategyName      `json:"strategy"`
	Currency  SupportedCurrency `json:"currency"`
}

type Strategy struct {
	Name         StrategyName      `json:"name"`
	Currency     SupportedCurrency `json:"currency"`
	OpenOffers   []Offer           `json:"open_offers"`
	ClosedOffers []Offer           `json:"closed_offers"`
}

type Offer struct {
	ClientOrderId         string                `json:"client_order_id"`
	ProductId             string                `json:"product_id"`
	Side                  Side                  `json:"side"`
	Config                OrderConfiguration    `json:"order_configuration"`
	SelfTradePreventionId SelfTradePreventionID `json:"self_trade_prevention_id"`
	RetailPortfolioId     string                `json:"retail_portfolio_id"`
}

type OrderConfiguration struct {
	LimitLimitGTD LimitLimitGTD `json:"limit_limit_gtd"`
}

type LimitLimitGTD struct {
	BaseSize   string `json:"base_size"`   // Amount of base currency to spend on order
	LimitPrice string `json:"limit_price"` // Ceiling price for which the order should get filled.
	EndTime    string `json:"end_time"`    // RFC3339 Timestamp
	PostOnly   bool   `json:"post_only"`   // If true, order should only make liquidity - maker commission charged.
}

type CoinbaseOrderPlacedResponse struct {
	Success         bool                               `json:"success"`
	FailureReason   string                             `json:"failure_reason"`
	OrderID         string                             `json:"order_id"`
	SuccessResponse CoinbaseOrderPlacedSuccessResponse `json:"success_response"`
	ErrorResponse   CoinbaseOrderPlacedErrorResponse   `json:"error_response"`
}

type CoinbaseOrderPlacedSuccessResponse struct {
	OrderID       string `json:"order_id"`
	ProductID     string `json:"product_id"`
	Side          string `json:"side"`
	ClientOrderID string `json:"client_order_id"`
}

type CoinbaseOrderPlacedErrorResponse struct {
	Error                 string `json:"error"`
	Message               string `json:"message"`
	ErrorDetails          string `json:"error_details"`
	PreviewFailureReason  string `json:"preview_failure_reason"`
	NewOrderFailureReason string `json:"new_order_failure_reason"`
}

type CoinbaseOrderPreviewResponse struct {
	OrderTotal       string   `json:"order_total"`
	CommissionTotal  string   `json:"commission_total"`
	Errors           []string `json:"errs"`
	Warnings         []string `json:"warning"`
	QuoteSize        string   `json:"quote_size"`
	BaseSize         string   `json:"base_size"`
	BestBid          string   `json:"best_bid"`
	BestAsk          string   `json:"best_ask"`
	IsMax            bool     `json:"is_max"`
	OrderMarginTotal string   `json:"order_margin_total"`
	Leverage         string   `json:"leverage"`
	ShortLeverage    string   `json:"short_leverage"`
	Slippage         string   `json:"slippage"`
}

type Side string

const (
	BUY  Side = "BUY"
	SELL Side = "SELL"
)

type StrategyName string

const (
	HODL StrategyName = "HODL"
)

type SupportedCurrency string

const (
	ETH SupportedCurrency = "ETH"
)

type SelfTradePreventionID string

const (
	Default SelfTradePreventionID = "ipopcorn-investment-manager"
)
