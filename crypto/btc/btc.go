package btc

import (
	"fmt"
	"time"

	"xox_tokens/configs"
	"xox_tokens/bus"

	"github.com/StorjPlatform/gocoin"
)

// Dial starts connection with BTCService via tcp.
func Dial(NATSPort uint) (err error) {
	bc, err := bus.Dial(fmt.Sprintf("http://0.0.0.0:%v", NATSPort))
	defer bc.Close()

	address := bus.SubjectAddress

	for {
		bc.Publish(address, &bus.EventAddress{Address: string(bus.SubjectAddress)})
		time.Sleep(time.Minute * 2)
	}
	return
}

func GenerateKey() (wif string, addr string, err error) {
	key, err := gocoin.GenerateKey(configs.GetConfigs().Crypto.BTCTestnet)
	if err != nil { return }

	wif = key.Priv.GetWIFAddress()
	addr, _ = key.Pub.GetAddress()
	return
}

