package types

type State struct {
	LastUpdated string      `json:"last_updated"`
	Portfolios  []Portfolio `json:"portfolios"`
}
