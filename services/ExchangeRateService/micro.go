package main

import (
	"context"
_	"fmt"
	"strconv"

	uErr "xox_tokens/errors"

	micro "github.com/micro/go-micro"
	proto "xox_tokens/services/ExchangeRateService/proto"
)

type Currency struct{}

func (c *Currency) GetRate(ctx context.Context, req *proto.GetRateReq, rsp *proto.GetRateResp) error {
	var currencyRatio float64
	switch req.BaseCurrency {
	case "BTC":
		currencyRatio = GetBTCRate()
	case "BTCETH":
		currencyRatio = GetBTCETHRate()
	case "FiatETH":
		currencyRatio = GetFiatETHRate()
	default:
		currencyRatio = GetETHRate()
	}
	rsp.CurrencyRatio = strconv.FormatFloat(currencyRatio, 'f', -1, 64)
	return nil
}

func microMain() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("Currency"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterCurrencyHandler(service.Server(), new(Currency))

	// Run the server
	if err := service.Run(); err != nil {
		uErr.Fatal(err)
		return
	}
}
