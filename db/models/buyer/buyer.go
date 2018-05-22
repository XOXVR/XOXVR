package buyer

import (
	"xox_tokens/db"
	"xox_tokens/db/types"
)

// New adds new buyer.
func New(ethAddr, btcAddr string) (err error) {
	var transactions []*types.Transaction

	err = db.Instance.Insert(&types.Buyer{
		EthAddr: ethAddr,
		BtcAddr: btcAddr,
		Transactions: transactions,
	})

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

	return &result, true, nil
}
