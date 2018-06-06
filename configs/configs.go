package configs

import (
	"math/big"

	"github.com/BurntSushi/toml"
)

type Configs struct {
	Common   Common
	Server   Server
	Crypto   Crypto
	Services Services
	DB       map[string]DB `toml:"database"`
	Stripe   map[string]Stripe
}

type Common struct {
	Dev        bool
	LogOutPath string
}

type Server struct {
	GRPCPort uint `toml:"gRPCPort"`
	RESTPort uint `toml:"restPort"`
	NATSPort uint `toml:"natsPort"`
}

type Crypto struct {
	// BTC
	BTCAddr          string `toml:"btcAddress"`
	BTCTestnet       bool   `toml:"btcTestnet"`

	// ETH
	ETHNetworkId     int `toml:"ethNetworkId"`
	InfuraToken      string

	// Contract
	CrowdsaleAddress string
	OwnerPrivateKey  string
	OwnerAddress     string
	GasPrice         big.Int
	GasLimit         uint64
}

type Services struct {
	ETHServicePort uint `toml:"ethServicePort"`
	BTCServicePort uint `toml:"btcServicePort"`
}

type DB struct {
	Host         string
	Port         uint
	Name         string
	User         string
	Password     string
}

type Stripe struct {
	PublishableKey 	 string
	SecretKey      	 string
}

var configs Configs

// ParseConfigs parses and returns configs from toml file.
func ParseConfigs(confPath string) (*Configs, error) {
	if _, err := toml.DecodeFile(confPath, &configs); err != nil {
		return nil, err
	}

	return &configs, nil
}

 //GetConfigs return configs.
 //You can get it from any place
 func GetConfigs() *Configs { return &configs }
