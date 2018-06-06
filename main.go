// Need to update .bashrc
// export MICRO_REGISTRY=mdns

package main

import (
	"flag"
	"os"
	"log"
	"context"
	"time"
	"strconv"
	"math/big"

	uErr "xox_tokens/errors"
	"xox_tokens/db"
	"xox_tokens/utils"
	"xox_tokens/configs"
	"xox_tokens/crypto/btc"
	"xox_tokens/crypto/eth"
	"xox_tokens/api/routers"
	"xox_tokens/services/ExchangeRateService/proto"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/cli"
)

func main() {
	prod := flag.Bool("prod", false, "Run in production mode.")
	flag.Parse()

	// Run watchers.
	defer utils.RecoverWatcher(shutdown)
	go utils.ShutdownWatcher(shutdown)

	// Parse configs.
	conf, err := configs.ParseConfigs("./configs.toml")
	if err != nil {
		uErr.Fatal(err, "failed to parse configs")
	}
	conf, err = configs.ParseConfigs("./configs_local.toml")
	if err != nil {
		uErr.Fatal(err, "failed to parse local configs")
	}

	// Setup prod env.
	envType := "dev"
	if *prod {
		envType = "prod"
		conf.Common.Dev = false

		err := utils.SetupLogFile(conf.Common.LogOutPath)
		if err != nil {
			uErr.Fatal(err, "failed to setup log file")
		}

		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("Start in dev mode")
	}

	// Initiate database.
	err = db.Initiate(conf.DB[envType])
	if err != nil {
		uErr.Fatal(err, "failed to initialize db")
	}
	defer db.Instance.Close()

	// Start connection with ETH provider.
	err = eth.Dial(&conf.Crypto)
	if err != nil {
		uErr.Fatal(err, "failed ETH Dial")
	}

	go eth.Refund()

	go func() {
		// Create a new service. Optionally include some options here.
		service := micro.NewService(
			micro.Name("Currency.client"),
			micro.Flags(
				cli.BoolFlag{
					Name:  "prod",
					Usage: "This is a prod flag",
				},
			),
		)
		service.Init(
		)

		// Create new greeter client
		greeter := proto.NewCurrencyService("Currency", service.Client())

		for {
			// Call the greeter
			rsp, err := greeter.GetRate(context.TODO(), &proto.GetRateReq{BaseCurrency: "FiatETH"})
			if err != nil {
				log.Println(err)
			}

			rate, err := strconv.ParseFloat(rsp.CurrencyRatio, 64)
			if err != nil {
				log.Println(err)
			}

			irate := int64(rate * float64(1e18))

			eth.SetRate(big.NewInt(irate))

			time.Sleep(time.Hour * 12)
		}
	}()

	// Template for REST
	go routers.RoutersStart(conf, envType)

	// Start connection with BTCService.
	go func() {
		log.Println("BTC service connection...")
		// TODO: Rename BTCAddr to BTCTrackedAddress
		if err = btc.Dial(conf.Server.NATSPort); err != nil {
			uErr.Fatal(err, "Failed to dial BTC Service")
		}
	}()

	select {}
}

func shutdown(fatal bool, r interface{}) {
	if fatal {
		os.Exit(1)
	}

	log.Println("Shutdown")
	os.Exit(0)
}
