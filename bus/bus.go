package bus

import (
	"github.com/nats-io/go-nats"
)

type Bus struct {
	conn *nats.EncodedConn

	// Subscriptions.
	subs map[Subject]*nats.Subscription
}

func Dial(addr string) (bus *Bus, err error) {
	nc, err := nats.Connect(addr)
	if err != nil { return }
	conn, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil { return }

	return &Bus{
		conn: conn,
		subs: map[Subject]*nats.Subscription{},
	}, nil
}

func (b *Bus) Close() {
	b.conn.Close()
}

func (b *Bus) Publish(subject Subject, msg interface{}) error {
	return b.conn.Publish(string(subject), msg)
}

func (b *Bus) Subscribe(subject Subject, cb nats.Handler) (err error) {
	sub, err := b.conn.Subscribe(string(subject), cb)
	if err != nil { return }

	b.subs[subject] = sub
	return
}

func (b *Bus) Unsubscribe(subject Subject) error {
	if sub, ok := b.subs[subject]; ok {
		delete(b.subs, subject)
		return sub.Unsubscribe()
	}

	return nil
}
