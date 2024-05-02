package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/server/util"
	"github.com/iPopcorn/investment-manager/types"
)

type HandleTransferFundsArgs struct {
	Client *infrastructure.InvestmentManagerExternalHttpClient
	Writer http.ResponseWriter
	Req    *http.Request
}

func HandleTransferFunds(args HandleTransferFundsArgs) {
	handlerName := "HandleTransferFunds: "

	if args.Req.Method != http.MethodPost {
		util.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Invalid http method, wanted %s got %s", http.MethodPost, args.Req.Method))

		return
	}

	body := args.Req.Body

	defer body.Close()

	bodyData, err := ioutil.ReadAll(body)

	if err != nil {
		log.Printf(handlerName+"Failed to read body from request: %v\n", err)
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	var reqBody types.TransferRequest

	err = json.Unmarshal(bodyData, &reqBody)

	if err != nil {
		log.Printf(handlerName + "Failed to deserialize request")
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	fundsToTransfer, err := strconv.ParseFloat(reqBody.Amount, 64)

	if err != nil {
		log.Printf(handlerName+"Invalid request\nCould not convert 'Amount' to float64\nGiven: %q\n%v", reqBody.Amount, err)
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	senderPortfolioDetails, err := util.PortfolioDetails(args.Client, reqBody.SenderID)

	if err != nil {
		log.Printf(handlerName+"Failed to get portfolio details for sender. Given: %q\n", reqBody.SenderID)
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	// get available funds, assume GBP for now.
	var senderAvailableFunds float64
	for _, position := range senderPortfolioDetails.Breakdown.SpotPositions {
		if position.Asset == "GBP" {
			senderAvailableFunds = position.AvailableToTradeFiat
			break
		}
	}

	if senderAvailableFunds < fundsToTransfer {
		log.Printf(handlerName+"Sender does not have enough funds to transfer\nAvailable funds: %f\nFunds to transfer %f\n", senderAvailableFunds, fundsToTransfer)

		args.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := util.TransferFunds(args.Client, &reqBody)

	if err != nil {
		log.Printf(handlerName+"Failed to transfer funds\nerror: %+v\nrequest: %+v", err, reqBody)
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	log.Printf(handlerName+"Transfer funds success!\nresp: %q\n", string(resp))
	util.WriteResponse(args.Writer, resp, nil)
}
