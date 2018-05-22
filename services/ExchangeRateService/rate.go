package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"xox_tokens/constants"
)

//ExchangeRateProvider describes service, which provides API for getting cryptocurrency to fiat currency exchange rate
type ExchangeRateProvider struct {
	Name   string
	ETHURL string
	BTCURL string
	USDURL string
}

//CryptonatorTicker is struct, which describes Cryptonator exchange ticker with information about currency price and 24h trading volume
type CryptonatorTicker struct {
	Base   string `json:"base"`
	Target string `json:"target"`
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Change string `json:"change"`
}

//CryptonatorResponse is a struct, which describes response from Cryptonator with exchange rate and
type CryptonatorResponse struct {
	Ticker    CryptonatorTicker `json:"ticker"`
	Timestamp int               `json:"timestamp"`
	Success   bool              `json:"success"`
	Error     string            `json:"error"`
}

type Rates struct {
	mux *sync.Mutex

	ETH float64
	BTC float64
	BTCETH float64
}

var (
	currentRates = Rates{&sync.Mutex{}, 0, 0, 0}

	RateProviders = [3]*ExchangeRateProvider{{
		"CryptoCompare",
		"https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=",
		"https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=",
		"https://min-api.cryptocompare.com/data/price?fsym=USD&tsyms=",
	}, {
		"Cryptonator",
		"https://api.cryptonator.com/api/ticker/eth-",
		"https://api.cryptonator.com/api/ticker/btc-",
		"https://api.cryptonator.com/api/ticker/usd-",
	}, {
		"CoinMarketCap",
		"https://api.coinmarketcap.com/v2/ticker/1027/?convert=",
		"https://api.coinmarketcap.com/v2/ticker/1/?convert=",
		"https://api.coinmarketcap.com/v2/ticker/1/?convert=", // NOTE: coinmarketcap api doesn't support fiat to coin conversion, need to use 1/rate
	}}
)

// GetAverageRate performs requests to CryptoCompare, Cryptonator and CoinMarketCap
// for (ETH|BTC)/symbol exchange rate and returns average rate.
func GetAverageRate(forCurrency uint8, symbol string, c *Config) float64 {
	rates := make(chan float64, 3)

	go func() {
		rates <- GetCryptoCompareRate(forCurrency, symbol, c)
	}()

	go func() {
		rates <- GetCryptonatorRate(forCurrency, symbol, c)
	}()

	go func() {
		rates <- GetCoinMarketCapRate(forCurrency, symbol, c)
	}()

	count := 0
	validCount := 0.0
	sum := 0.0
	select {
	case val := <-rates:
		count++
		if val != -math.MaxFloat64 && val > 0 {
			sum += val
			validCount++
		}
		if count == 3 {
			break
		}
	}

	if validCount == 0 {
		return -math.MaxFloat64
	}

	return sum / validCount
}

func SetRate(forCurrency uint8, rate float64) {
	currentRates.mux.Lock()
	defer currentRates.mux.Unlock()

	switch forCurrency {
		case constants.BTC:
			currentRates.BTC = rate
		case constants.ETH:
			currentRates.ETH = rate
		case constants.BTCETH:
			currentRates.BTCETH = rate
	}
}

func GetETHRate() float64 {
	return currentRates.ETH
}

func GetBTCRate() float64 {
	return currentRates.BTC
}

func GetFiatETHRate() float64 {
	return 1 / currentRates.ETH
}

func GetBTCETHRate() float64 {
	return currentRates.BTCETH
}

//GetCryptoCompareRate performs GET request to CryptoCompare and returns ETH/symbol exchange rate
func GetCryptoCompareRate(forCurrency uint8, symbol string, conf *Config) float64 {
	rate := -math.MaxFloat64
	if conf == nil {
		return rate
	}

	var fullUrl string
	if forCurrency == constants.ETH {
		fullUrl = RateProviders[0].ETHURL + symbol
	} else {
		fullUrl = RateProviders[0].BTCURL + symbol
	}

	conf.Logger.Println("Querying CryptoCompare for exchange rate...")

	// Perform c.Retries of request retries
	for retries := 0; retries < conf.Retries; retries++ {
		resp, err := http.Get(fullUrl)
		if err != nil {
			conf.Logger.Println("CryptoCompare request error: " + err.Error())
			continue
		}
		dec := json.NewDecoder(resp.Body)
		// Manually read JSON response, looking for key equal to fiat symbol
		for {
			t, tErr := dec.Token()
			if tErr == io.EOF {
				break
			}
			if tErr != nil {
				fmt.Println("JSON response decoding error: " + tErr.Error())
				continue
			}
			//if we have found fiat symbol in response, than next token will be exchange rate
			if fmt.Sprintf("%v", t) == symbol {
				val, _ := dec.Token()
				rate = val.(float64)
				break
			}
		}
		resp.Body.Close()
		//stop retrying requests if we have found some non-default value
		if rate != -math.MaxFloat64 {
			break
		}
	}

	return rate
}

//GetCryptonatorRate performs GET request to Cryptonator and returns ETH/symbol exchange rate
func GetCryptonatorRate(forCurrency uint8, symbol string, conf *Config) float64 {
	rate := -math.MaxFloat64
	if conf == nil {
		return rate
	}

	var fullUrl string
	if forCurrency == constants.ETH {
		fullUrl = RateProviders[1].ETHURL + strings.ToLower(symbol)
	} else {
		fullUrl = RateProviders[1].BTCURL + strings.ToLower(symbol)
	}

	conf.Logger.Println("Querying Cryptonator for exchange rate...")

	// Perform conf.Retries of request retries
	for retries := 0; retries < conf.Retries; retries++ {
		resp, err := http.Get(fullUrl)
		if err != nil {
			conf.Logger.Println("CryptoCompare request error: " + err.Error())
			continue
		}
		response := &CryptonatorResponse{}
		err = json.NewDecoder(resp.Body).Decode(response)
		resp.Body.Close()
		if err != nil {
			conf.Logger.Println("Failed to decode cryptonator response")
			continue
		}

		if response.Error != "" {
			conf.Logger.Println("Cryptonator error: ", response.Error)
			continue
		}

		r, err := strconv.ParseFloat(response.Ticker.Price, 64)
		if err != nil {
			conf.Logger.Println("Failed to parse float in response")
			continue
		}

		rate = r

		//stop retrying requests if we have found some non-default value
		if rate != -math.MaxFloat64 {
			break
		}
	}

	return rate
}

//GetCoinMarketCapRate performs GET request to CoinMarketCap and returns ETH/symbol exchange rate
func GetCoinMarketCapRate(forCurrency uint8, symbol string, conf *Config) float64 {
	rate := -math.MaxFloat64
	if conf == nil {
		return rate
	}

	var fullUrl string
	if forCurrency == constants.ETH {
		fullUrl = RateProviders[2].ETHURL + symbol
	} else {
		fullUrl = RateProviders[2].BTCURL + symbol
	}

	responseField := "price_" + strings.ToLower(symbol)

	conf.Logger.Println("Querying CoinMarketCap for exchange rate...")

	// Perform c.Retries of request retries
	for retries := 0; retries < conf.Retries; retries++ {
		resp, err := http.Get(fullUrl)
		if err != nil {
			conf.Logger.Println("CryptoCompare request error: " + err.Error())
			continue
		}
		dec := json.NewDecoder(resp.Body)
		// Manually read JSON response, looking for key equal to fiat symbol
		for {
			t, tErr := dec.Token()
			if tErr == io.EOF {
				break
			}
			if tErr != nil {
				fmt.Println("JSON response decoding error: " + tErr.Error())
				continue
			}
			//if we have found fiat symbol in response, than next token will be exchange rate
			if fmt.Sprintf("%v", t) == responseField {
				val, _ := dec.Token()
				strRate := fmt.Sprintf("%v", val)
				r, err := strconv.ParseFloat(strRate, 64)
				if err == nil {
					rate = r
				}
				break
			}
		}
		resp.Body.Close()
		//stop retrying requests if we have found some non-default value
		if rate != -math.MaxFloat64 {
			break
		}
	}

	return rate
}
