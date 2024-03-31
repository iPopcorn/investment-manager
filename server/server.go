package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/iPopcorn/investment-manager/infrastructure"
)

type InvestmentManagerHTTPServer struct {
	client infrastructure.InvestmentManagerExternalHttpClient
}

func (s *InvestmentManagerHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %v\n", r)

	route, args := getRouteAndArgsFromPath(r.URL.Path)

	switch route {
	case "portfolios":
		s.handlePortfolio(w, r, args)
		return
	default:
		log.Printf("Route not found: %q\n", route)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func GetInvestmentManagerHTTPServer() *InvestmentManagerHTTPServer {
	httpClient := infrastructure.GetInvestmentManagerExternalHttpClient()

	return &InvestmentManagerHTTPServer{
		client: *httpClient,
	}
}

func (s *InvestmentManagerHTTPServer) handlePortfolio(w http.ResponseWriter, r *http.Request, args []string) {
	w.Header().Set("Content-Type", "application/json")
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"

	if r.Method == http.MethodPost {
		body := r.Body

		defer body.Close()

		bodyData, err := ioutil.ReadAll(body)

		if err != nil {
			log.Printf("Failed to read body from request: %v\n", err)
			writeResponse(w, nil, err)
		}

		resp, err := s.client.Post(url, bodyData)

		writeResponse(w, resp, err)
	} else {
		if len(args) == 1 {
			portfolioUUID := args[0]
			url = url + "/" + portfolioUUID
			resp, err := s.client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolio details from URL: %q\nError: %v", url, err)
			}

			writeResponse(w, resp, err)
		} else {
			resp, err := s.client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolios from URL: %q\nError: %v", url, err)
			}

			writeResponse(w, resp, err)
		}
	}

	log.Printf("Request handled successfully!")
}

func writeResponse(w http.ResponseWriter, response []byte, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)

	if err != nil {
		log.Println("Failed to write response to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getRouteAndArgsFromPath(path string) (string, []string) {
	rawPath := strings.TrimPrefix(path, "/")
	log.Printf("rawPath: %v", rawPath)
	pathTokens := strings.Split(rawPath, "/")
	log.Printf("pathTokens: %v", pathTokens)
	route := pathTokens[0]
	log.Printf("Route: %q\n", route)
	args := []string{}

	for i := 1; i < len(pathTokens); i++ {
		args = append(args, pathTokens[i])
	}

	return route, args
}
