package transaction

import (
	"xox_tokens/db"
	"xox_tokens/db/types"
)

// FindByHash returns transaction with particular hash.
func FindByHash(hash string) (*types.Transaction, bool, error) {
	var result types.Transaction
	err := db.Instance.Model(&result).
		Where("hash = ?", hash).
		Select()
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &result, true, nil
}
