package events

import (
	"xox_tokens/bus"
)

// Subscribe subscribes on events.
func Subscribe(b *bus.Bus) (err error) {
	err = b.Subscribe(bus.SubjectAddress, onNewAddress)
	if err != nil { return }
	err = b.Subscribe(bus.SubjectValue, onNewValue)
	return
}
