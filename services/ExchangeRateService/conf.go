package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type Config struct {
	Retries       int
	NetworkId     int
	GasLimit      uint64
	InfuraToken   string
	FiatSymbol    string
	EthConnection *ethclient.Client
	Transactor    *bind.TransactOpts
	Logger        *log.Logger
	Contract      common.Address
	UpdateRate    time.Duration
}

var conf *Config

func GetConfig() *Config {
	if conf == nil {
		conf = loadConfig()
	}
	return conf
}

func loadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("rate-conf")
	viper.SetConfigType("yaml")
	viper.SetDefault("fiat_symbol", "USD")
	viper.SetDefault("network_id", 3) // default value for network is Rinkeby
	viper.SetDefault("retries", 1)    // default value for number of retries is 1
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read config file")
	}

	c := &Config{}

	c.InfuraToken = viper.GetString("infura_token")
	c.NetworkId = viper.GetInt("network_id")

	contractAddress := viper.GetString("contract_address")
	if !common.IsHexAddress(common.HexToAddress(contractAddress).Hex()) {
		fmt.Println("Valid contract address is required")
		return nil
	}
	c.Contract = common.HexToAddress(contractAddress)

	c.UpdateRate = viper.GetDuration("update_rate")

	privateKey := strings.TrimPrefix(viper.GetString("private_key"), "0x")
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println("Failed to parse secp256k1 key for transaction signing")
		return nil
	}
	c.Transactor = bind.NewKeyedTransactor(key)

	c.FiatSymbol = viper.GetString("fiat_symbol")
	c.Retries = viper.GetInt("retries")

	c.GasLimit = uint64(viper.GetInt64("gas_limit"))

	return c
}
