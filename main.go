// Need to update .bashrc
// export MICRO_REGISTRY=mdns

package main

import (
	"flag"
	"os"
	"log"
	"fmt"
	"context"
_	"time"
	"strconv"
	"math/big"

	uErr "xox_tokens/errors"
	"xox_tokens/db"
	"xox_tokens/utils"
	"xox_tokens/configs"
_	"xox_tokens/crypto"
_	"xox_tokens/crypto/btc"
	"xox_tokens/crypto/eth"
_	"xox_tokens/service"

	micro "github.com/micro/go-micro"
	proto "xox_tokens/services/ExchangeRateService/proto"
)

func main() {
	prod := flag.Bool("prod", false, "Run in production mode.") // TODO: Uncomment this line
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
	}

	fmt.Println("==> " + envType) // TODO: Delete this line

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
		service := micro.NewService(micro.Name("greeter.client"))
		service.Init()

		// Create new greeter client
		greeter := proto.NewCurrencyService("Currency", service.Client())

		// Call the greeter
		rsp, err := greeter.GetRate(context.TODO(), &proto.GetRateReq{BaseCurrency: "FiatETH"})
		if err != nil {
			fmt.Println(err)
		}

		rate, err := strconv.ParseFloat(rsp.CurrencyRatio, 64)
		if err != nil {
			fmt.Println(err)
		}

		irate := int64(rate * float64(1e18))

		eth.SetRate(big.NewInt(irate))

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
