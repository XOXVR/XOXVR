package errors

import (
	"log"
	"fmt"
	"errors"
)

// Error codes.
const (
	ErrInsertToDBCode   = iota + 1
	ErrSelectFromDBCode
)

// TODO: Replace Error -> Err
// Error messages.
const (
	ErrSelectFromDB       = "failed to select from DB"
	ErrInsertToDB         = "failed to insert to DB"
	ErrFindAddress        = "address not found"
	ErrConnectBTCService  = "btc service shutdown"
	ErrConnectRateService = "rate service shutdown"

	// API
	ErrValidateBTCAddress = "BTC address is not valid"
	ErrValidateETHAddress = "ETH address is not valid"


	ErrTXTimedOut         = "transaction timed out"
	ErrReceiptStatus      = "receipt.Status is 0"
	ErrUpdateGasLimit     = "failed to update gas limit"
	ErrOwner              = "permission denied, please check ownerAddress and ownerPrivateKey"
	ErrICONotStarted      = "ico not started"
	ErrICOFinished        = "ico finished"

	ErrUnknown            = "unknown error"
)

// LogError records error in log.
func LogError(err error, description ...interface{}) {
	log.Println("ERROR -", fmt.Sprint(description...), ":", getErrStr(err))
}

// Fatal records fatal error in log and stops process.
func Fatal(err error, description ...interface{}) {
	panic(fmt.Sprint(fmt.Sprint(description...), " : ", getErrStr(err)))
}

// Combine creates and returns error with description.
func Combine(err error, description ...interface{}) error {
	return errors.New(fmt.Sprintf("%s: %s", fmt.Sprint(description...), getErrStr(err)))
}

// ErrCode
func CombineErrCode(errCode int) (error) {
	return errors.New(fmt.Sprintf("code %d", errCode))
}

func getErrStr(err error) (errStr string) {
	if err != nil { errStr = err.Error() }
	return
}
