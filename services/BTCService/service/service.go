package service

import (
	"fmt"

	"xox_tokens/services/BTCService/btc"
	"xox_tokens/bus"
	"xox_tokens/services/events"
)

// StartTCPServer starts TCP listening on certain port.
func StartTCPServer(NATSPort uint, _btcWatcher *btc.Watcher) (err error) {
	bc, err := bus.Dial(fmt.Sprintf("http://0.0.0.0:%v", NATSPort))
	defer bc.Close()

	events.Subscribe(bc)

	select {}

	return
}
