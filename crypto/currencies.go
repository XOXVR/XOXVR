package crypto

import (
	"net"
	"fmt"
	"bufio"
	"io"
	"bytes"
	"encoding/gob"

	uErr "xox_tokens/errors"
	"xox_tokens/types"
	"xox_tokens/constants"
	"sync"
)

type Rates struct {
	mux *sync.Mutex

	ETH float64
	BTC float64
}

var (
	rates = Rates{&sync.Mutex{}, 0, 0}
)

// Dial starts connection with ExchangeRateService via tcp.
func Dial(port uint) (err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil { return }

	response := bufio.NewReader(conn)
	for {
		message, err := response.ReadBytes(byte('\n'))
		switch err {
		case nil:
			go handleMessage(message)
		case io.EOF:
			return uErr.Combine(nil, uErr.ErrConnectRateService)
		default:
			return err
		}
	}

	return
}

func handleMessage(line []byte) {
	var message types.RateServiceResp

	r := bytes.NewReader(line)
	dec := gob.NewDecoder(r)
	err := dec.Decode(&message)
	if err != nil { return }

	if message.FiatCurrency != "EUR" {
		uErr.Fatal(nil, "actually supported EUR only")
	}

	rates.mux.Lock()
	switch message.Currency {
	case constants.ETH:
		rates.ETH = message.Value
		fmt.Println("ETH rate updated:", message.Value)
	case constants.BTC:
		rates.BTC = message.Value
		fmt.Println("BTC rate updated:", message.Value)
	default:
		break
	}
	rates.mux.Unlock()
}

// TODO: Might to use mutex here.
// GetETHRate returns current ETH rate.
func GetETHRate() float64 {
	return rates.ETH
}

// GetBTCRate returns current BTC rate.
func GetBTCRate() float64 {
	return rates.BTC
}
