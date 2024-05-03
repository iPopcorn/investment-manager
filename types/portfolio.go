package types

type TransferRequest struct {
	SenderID   string
	ReceiverID string
	Amount     string
}

type PortfolioResponse struct {
	Portfolios []Portfolio `json:"portfolios"`
}

type PortfolioCreatedResponse struct {
	Portfolio Portfolio `json:"portfolio"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Portfolio struct {
	Name               string      `json:"name"`
	Uuid               string      `json:"uuid"`
	Type               string      `json:"type"`
	Deleted            bool        `json:"deleted"`
	CurrentStrategy    *Strategy   `json:"current_strategy"`
	PreviousStrategies *[]Strategy `json:"previous_strategies"`
}

type PortfolioDetailsResponse struct {
	Breakdown Breakdown `json:"breakdown"`
}

type Breakdown struct {
	Portfolio         Portfolio         `json:"portfolio"`
	PortfolioBalances PortfolioBalances `json:"portfolio_balances"`
	SpotPositions     []SpotPositions   `json:"spot_positions"`
	// PerpPositions     []PerpPositions    `json:"perp_positions"`    // TODO - implement
	// FuturesPositions  []FuturesPositions `json:"futures_positions"` // TODO - implement
}

type SpotPositions struct {
	Asset                string  `json:"asset"`
	AccountUuid          string  `json:"account_uuid"`
	TotalBalanceFiat     float64 `json:"total_balance_fiat"`
	TotalBalanceCrypto   float64 `json:"total_balance_crypto"`
	AvailableToTradeFiat float64 `json:"available_to_trade_fiat"`
	Allocation           float64 `json:"allocation"`
	OneDayChange         float64 `json:"one_day_change"`
	CostBasis            Balance `json:"cost_basis"`
	AssetImgUrl          string  `json:"asset_img_url"`
	IsCash               bool    `json:"is_cash"`
}

type PortfolioBalances struct {
	TotalBalance               Balance `json:"total_balance"`
	TotalFuturesBalance        Balance `json:"total_futures_balance"`
	TotalCashEquivalentBalance Balance `json:"total_cash_equivalent_balance"`
	TotalCryptoBalance         Balance `json:"total_crypto_balance"`
	FuturesUnrealizedPnl       Balance `json:"futures_unrealized_pnl"`
	PerpUnrealizedPnl          Balance `json:"perp_unrealized_pnl"`
}

type Balance struct {
	Value    string
	Currency string
}
