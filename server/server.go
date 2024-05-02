package server

import (
	"log"
	"net/http"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/server/handlers"
	"github.com/iPopcorn/investment-manager/server/state"
	"github.com/iPopcorn/investment-manager/server/util"
	"github.com/iPopcorn/investment-manager/types"
)

type InvestmentManagerHTTPServer struct {
	client          infrastructure.InvestmentManagerExternalHttpClient
	stateRepository *state.StateRepository
	channels        []chan bool
}

type InvestmentManagerHTTPServerArgs struct {
	HttpClient      *infrastructure.InvestmentManagerExternalHttpClient
	StateRepository *state.StateRepository
	Channels        []chan bool
}

func GetDefaultInvestmentManagerHTTPServer() *InvestmentManagerHTTPServer {
	httpClient := infrastructure.GetInvestmentManagerExternalHttpClient()
	stateRepo := state.StateRepositoryFactory("")

	return &InvestmentManagerHTTPServer{
		client:          *httpClient,
		stateRepository: stateRepo,
	}
}

func InvestmentManagerHttpServerFactory(args InvestmentManagerHTTPServerArgs) *InvestmentManagerHTTPServer {
	return &InvestmentManagerHTTPServer{
		client:          *args.HttpClient,
		stateRepository: args.StateRepository,
		channels:        args.Channels,
	}
}

func (s *InvestmentManagerHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %v\n", r)

	route, args := util.GetRouteAndArgsFromPath(r.URL.Path)

	switch route {

	case string(types.Portfolios):
		handlePortfolioArgs := handlers.HandlePortfolioArgs{
			Client: &s.client,
			Writer: w,
			Req:    r,
			Args:   args,
		}

		handlers.HandlePortfolio(handlePortfolioArgs)
		return

	case string(types.ExecuteStrategy):
		executeStrategyArgs := handlers.HandleExecuteStrategyArgs{
			Client:          &s.client,
			Writer:          w,
			Req:             r,
			Args:            args,
			Channels:        s.channels,
			StateRepository: s.stateRepository,
		}

		handlers.HandleExecuteStrategy(executeStrategyArgs)
		return

	case string(types.TransferFunds):
		handleTransferFundsArgs := handlers.HandleTransferFundsArgs{
			Client: &s.client,
			Writer: w,
			Req:    r,
		}

		handlers.HandleTransferFunds(handleTransferFundsArgs)
		return

	default:
		log.Printf("Route not found: %q\n", route)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
