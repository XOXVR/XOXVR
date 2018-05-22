package utils

import (
	"log"
	"os"
	"os/signal"
	"math/big"
)

var (
	percents = new(big.Int).SetUint64(10)
	oneH = new(big.Int).SetUint64(100)
	gOne = new(big.Int).SetUint64(1)
)

// RecoverWatcher waiting for recover.
func RecoverWatcher(shutdownFunc func(bool, interface{})) {
	if r := recover(); r != nil {
		log.Println("----------------------------------")
		log.Println("SESSION END WITH FATAL ERROR -", r)
		log.Println("----------------------------------")
		shutdownFunc(true, r)
	}
}

// ShutdownWatcher waiting for os interrupt signal
// and handle it.
func ShutdownWatcher(shutdownFunc func(bool, interface{})) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	select {
	case <-sig:
		shutdownFunc(false, nil)
	}
}

// DoNTimeBeforeComplete calls function for N times
// while it returns error.
func DoNTimeBeforeComplete(count int, rFunc func(i int) error) (err error) {
	for i := 0; i <= count; i++ {
		err = rFunc(i)
		if err != nil {
			continue
		} else {
			return nil
		}
	}

	return
}

// IncreaseGasPrice increases gas price value on 10% + 1.
func IncreaseGasPrice(gasPrice *big.Int) {
	bPercents := new(big.Int).Set(gasPrice)
	bPercents.Mul(bPercents, percents)
	bPercents.Quo(bPercents, oneH)
	gasPrice.Add(gasPrice, bPercents)
	gasPrice.Add(gasPrice, gOne)
}
