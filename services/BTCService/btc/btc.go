package btc

import (
	"errors"
	"log"
	"time"
	"sync"
	"fmt"

	uErr "xox_tokens/errors"
	"xox_tokens/services/BTCService/store"
	"xox_tokens/services/BTCService/configs"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Watcher struct {
	mux              sync.Mutex
	client           *rpcclient.Client

	ActiveWatcherId  uint8

	OnNewValue       func(float64, string, string)
}

const (
	ACCOUNT        = "mainAccount"
	LAST_BLOCK_KEY = "lastHandledBTCBlock"
)

// StartRPCConnection starts connection to node.
func StartRPCConnection(conf configs.Bitcoin) (watcher *Watcher, err error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.RPCHost,
		User:         conf.User,
		Pass:         conf.Password,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil { return }

	watcher = &Watcher{sync.Mutex{}, client,0, nil}

	return
}

// IncreaseActiveWatcherId increases ActiveWatcherId.
func (w *Watcher) IncreaseActiveWatcherId() uint8 {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.ActiveWatcherId == 255 {
		w.ActiveWatcherId = 0
	} else {
		w.ActiveWatcherId++
	}
	return w.ActiveWatcherId
}

// CheckActiveWatcherIdChanged checks if ActiveWatcherId was changed.
func (w *Watcher) CheckActiveWatcherIdChanged(ownerId uint8) bool {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.ActiveWatcherId == ownerId { return true }
	return false
}

// StartWatchingAddress handles error from WatchAddress.
func (w *Watcher) StartWatchingAddress(addr string) {
	log.Println("Start watching address " + addr + "...")
	err := w.WatchAddress(addr)
	if err != nil {
		uErr.Fatal(err, "failed to watch address")
	}
}

// WatchAddress adds listener for address.
func (w *Watcher) WatchAddress(addr string) error {
	if addr == "" {
		return errors.New(uErr.ErrFindAddress)
	}

	var lastHandledBlock int64
	lastHandledBlockFromCache, found := store.Get(LAST_BLOCK_KEY)
	if found {
		lastHandledBlock = lastHandledBlockFromCache.(int64)
	}

	selfId := w.IncreaseActiveWatcherId()

	for {
		stillOwner := w.CheckActiveWatcherIdChanged(selfId)
		if !stillOwner { return nil }

		bCount, err := w.client.GetBlockCount()
		if err != nil { return err }
		log.Println("blocks count:", bCount)
		log.Println("last handled block:", lastHandledBlock)

		// If lastHandledBlock backward in time.
		backwardInTime := false
		if lastHandledBlock > 0 && bCount - lastHandledBlock > 1 {
			bCount = lastHandledBlock + 1
			backwardInTime = true
		} else if bCount == lastHandledBlock {
			sleep()
			continue
		}

		lastHash, err := w.client.GetBlockHash(bCount)
		if err != nil { return err }
		//log.Println("last block", lastHash)

		lastBlock, err := w.client.GetBlockVerbose(lastHash)
		if err != nil { return err }

		// Filter transactions.
		for _, txStr := range lastBlock.Tx {
			txHash, _ := chainhash.NewHashFromStr(txStr)

			// Get transaction.
			tx, err := w.client.GetRawTransactionVerbose(txHash)
			if err != nil { return err }

			// Filter addresses.
			for _, vout := range tx.Vout {
				for _, fAddr := range vout.ScriptPubKey.Addresses {
					//log.Println("Address:", fAddr)
					if fAddr == addr { // TODO: Remove redundant logs
						log.Println("------------------")
						log.Println("Transaction for", addr)
						log.Println("txStr:", txStr)
						log.Println("tx:", tx)
						log.Println("vout:", vout)
						log.Println("------------------")

						// Detect TX owner.
						in := tx.Vin[0]
						ownerTxHash, _ := chainhash.NewHashFromStr(in.Txid)
						ownerTxOutIndex := in.Vout

						ownerTx, err := w.client.GetRawTransactionVerbose(ownerTxHash)
						if err != nil {
							return uErr.Combine(err, "failed to get owner tx")
						}

						ownerOut := ownerTx.Vout[ownerTxOutIndex]
						ownerAddr := ownerOut.ScriptPubKey.Addresses[0]
						if w.OnNewValue != nil {
							w.OnNewValue(vout.Value, ownerAddr, tx.Hash)
						}

						log.Println("tx owner:", ownerAddr, "tx value", vout.Value)
					}
				}
			}
		}

		// Save last handled block.
		store.Set(LAST_BLOCK_KEY, bCount)
		lastHandledBlock = bCount

		// Sleep for 1 minute.
		if !backwardInTime {
			sleep()
		}
	}

	return nil
}

func sleep() {
	fmt.Println("sleeping for 2 minutes...")
	time.Sleep(2 * time.Minute)
}
