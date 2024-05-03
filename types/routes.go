package types

type Route string

const (
	Portfolios      Route = "portfolios"
	ExecuteStrategy Route = "execute-strategy"
	TransferFunds   Route = "transfer-funds"
)
