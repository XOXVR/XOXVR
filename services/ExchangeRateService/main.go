package main

import (
	"log"
	"math"
	"os"
	"time"
	"fmt"

	"xox_tokens/constants"
)

func main() {
	fmt.Println("==============================================") // TODO: Delete this line
	conf := GetConfig()
	if conf == nil {
		panic("Failed to read config file")
	}
	conf.Logger = log.New(os.Stdout, "ExchangeRateService: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	// Run rate watchers.
	go watchETHRate(conf)
	go watchBTCRate(conf)
	go watchBTCETHRate(conf)

	microMain()
}

func watchETHRate(conf *Config) {
	for {
		conf.Logger.Println("Updating ETH rate...")
		rate := GetAverageRate(constants.ETH, conf.FiatSymbol, conf)
		if rate == -math.MaxFloat64 {
			conf.Logger.Println("Failed to get ETH rate from 3 sources, skipping")
			time.Sleep(conf.UpdateRate * time.Minute)
			continue
		}

		conf.Logger.Println("Average ETH/" + conf.FiatSymbol + " exchange rate -", rate)

		SetRate(constants.ETH, rate)

		conf.Logger.Println("Sleeping for", int64(conf.UpdateRate), "minutes..")
		time.Sleep(conf.UpdateRate * time.Minute)
	}
}

func watchBTCRate(conf *Config) {
	for {
		conf.Logger.Println("Updating BTC rate...")
		rate := GetAverageRate(constants.BTC, conf.FiatSymbol, conf)
		if rate == -math.MaxFloat64 {
			conf.Logger.Println("Failed to get BTC rate from 3 sources, skipping")
			time.Sleep(conf.UpdateRate * time.Minute)
			continue
		}

		conf.Logger.Println("Average BTC/" + conf.FiatSymbol + " exchange rate -", rate)

		SetRate(constants.BTC, rate)

		conf.Logger.Println("Sleeping for", int64(conf.UpdateRate), "minutes..")
		time.Sleep(conf.UpdateRate * time.Minute)
	}
}

func watchBTCETHRate(conf *Config) {
	for {
		conf.Logger.Println("Updating USD to ETH rate...")
		rate := GetAverageRate(constants.BTC, "ETH", conf)
		if rate == -math.MaxFloat64 {
			conf.Logger.Println("Failed to get USD to ETH rate from 3 sources, skipping")
			time.Sleep(conf.UpdateRate * time.Minute)
			continue
		}

		conf.Logger.Println("Average BTC/ETH exchange rate -", rate)

		SetRate(constants.BTCETH, rate)

		conf.Logger.Println("Sleeping for", int64(conf.UpdateRate), "minutes..")
		time.Sleep(conf.UpdateRate * time.Minute)
	}
}
