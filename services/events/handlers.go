package events

import (
	"fmt"
	"xox_tokens/bus"
)

func onNewAddress(address *bus.EventAddress) {
	// TODO: Calculate values for game and user.
	fmt.Println("New address received:", address)
}

func onNewValue(value *bus.EventValue) {
	// TODO: Calculate values for game and user.
	fmt.Println("New value received:", value)
}
