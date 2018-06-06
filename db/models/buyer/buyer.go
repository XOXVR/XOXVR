package buyer

import (
	"encoding/base64"
	"errors"

	"xox_tokens/db"
	"xox_tokens/db/types"
)

// New adds new buyer.
func New(ethAddr, btcAddr, btcWif string) (buyer *types.Buyer, err error) {
	var transactions []*types.Transaction
	encodedWif := base64.StdEncoding.EncodeToString([]byte(btcWif))

	err = db.Instance.Insert(&types.Buyer{
		EthAddr: ethAddr,
		BtcAddr: btcAddr,
		BtcAddrWif: encodedWif,
		Transactions: transactions,
	})

	if err != nil {
		return nil, err
	}

	buyer, _, err = FindByETHAddress(ethAddr)

	return
}

// FindByBTCAddress returns buyer by ethAddr.
func FindByBTCAddress(btcAddr string) (*types.Buyer, bool, error) {
	var result types.Buyer
	err := db.Instance.Model(&result).
		Where("btc_addr = ?", btcAddr).
		Select()
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	decodedWif, err := base64.StdEncoding.DecodeString(result.BtcAddrWif)
	if err != nil {
		return nil, false, errors.New("WIF decode error")
	}
	result.BtcAddrWif = string(decodedWif)

	return &result, true, nil
}

// FindByETHAddress returns buyer by ethAddr.
func FindByETHAddress(ethAddr string) (*types.Buyer, bool, error) {
	var result types.Buyer
	err := db.Instance.Model(&result).
		Where("eth_addr = ?", ethAddr).
		Select()
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	decodedWif, err := base64.StdEncoding.DecodeString(result.BtcAddrWif)
	if err != nil {
		return nil, false, errors.New("WIF decode error")
	}
	result.BtcAddrWif = string(decodedWif)

	return &result, true, nil
}
