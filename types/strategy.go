package types

type ExecuteStrategyRequest struct {
	Portfolio string `json:"portfolio"`
	Strategy  string `json:"strategy"`
	Currency  string `json:"currency"`
}
