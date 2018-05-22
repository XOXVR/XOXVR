package eth

import (
	"fmt"
	"log"
	"errors"
	"math/big"
	"time"
_	"math"

	uErr "xox_tokens/errors"
//	uCrypto "xox_tokens/crypto"
	"xox_tokens/gocontracts"
	"xox_tokens/configs"
	"xox_tokens/utils"

	"golang.org/x/net/context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
_	"github.com/ethereum/go-ethereum/core"
)

var (
	client     *ethclient.Client

	auth       *bind.TransactOpts
	dAuth      *bind.TransactOpts
	session    *gocontracts.CrowdsaleSession
	euroCents  *big.Int
)

// Dial connects to ETH provider create sessions for contracts.
func Dial(conf *configs.Crypto) (err error) {
	providerURL := utils.GetInfuraProviderUrl(conf.ETHNetworkId, conf.InfuraToken)

	client, _ = ethclient.Dial(providerURL)
	if err != nil { return }

	// Create TransactOpts.
	auth, err = createAuth(conf.OwnerPrivateKey, &conf.GasPrice, conf.GasLimit)
	if err != nil {
		return uErr.Combine(err, "failed to create auth")
	}

	// Create Crowdsale session.
	session, err = createCrowdsaleSession(conf.CrowdsaleAddress)
	if err != nil {
		return uErr.Combine(err, "failed to create Crowdsale session")
	}

	return nil
}

// TODO: Edit this func
func Refund() error {
	ico, err := session.ICO()
	softCap, err := session.SOFTCAP()

	if err != nil {
		uErr.Fatal(err)
		return err
	}

	if sleepTime := ico.End.Int64() - time.Now().Unix(); sleepTime > 0  {
		time.Sleep(time.Duration(sleepTime) + time.Hour)
	}

	salesTokens, err := session.SalesTokens()

	if err != nil {
		uErr.Fatal(err)
		return err
	}

	if ico.End.Int64() < time.Now().Unix() && salesTokens.Int64() < softCap.Int64() {

		numberInvestors, err := session.ReturnNumberInvestors()
		if err != nil {
			uErr.Fatal(err)
			return err
		}

		countRefund, err := session.CountRefund()
		for ; countRefund.Int64() < numberInvestors.Int64(); {
			if err != nil {
				uErr.Fatal(err)
				return err
			}

			tx, err := session.Refund()
			if err != nil {
				uErr.Fatal(err)
				return err
			}

			_, err = getReceipt(tx, true)
			if err != nil {
				uErr.Fatal(err)
				return err
			}
			log.Println("Refund iteration complete")
			countRefund, err = session.CountRefund()
		}
	}

	log.Println("Refund complete")
	return nil
}

func SetRate(rate *big.Int) error {
	tx, err := session.SetRate(rate)
	if err != nil {
		uErr.Fatal(err)
		return err
	}

	_, err = getReceipt(tx, true)
	if err != nil {
		uErr.Fatal(err)
		return err
	}
	log.Println("tx complete")

	return nil
}

func getGasLimit() (uint64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return block.GasLimit(), nil
}

// UpdateGasLimit asks for new gas limit value
// and sets it for TransactOpts.
func UpdateGasLimit() (err error) {
	gasLimit, err := getGasLimit()
	if err != nil { return }

	session.TransactOpts.GasLimit = gasLimit
	return
}

func createAuth(key string, gasPrice *big.Int, gasLimit uint64) (auth *bind.TransactOpts, err error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil { return }

	log.Println("GasPrice is:", gasPrice, "GasLimit is:", gasLimit)

	auth = bind.NewKeyedTransactor(privateKey)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	return
}

func createCrowdsaleSession(contractAddr string) (*gocontracts.CrowdsaleSession, error) {
	addr := common.HexToAddress(contractAddr)
	contract, err := gocontracts.NewCrowdsale(addr, client)
	if err != nil { return nil, err }

	session := &gocontracts.CrowdsaleSession{
		Contract: contract,
		CallOpts: bind.CallOpts{
			Pending: true,
			From:    auth.From,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasPrice: auth.GasPrice,
			GasLimit: auth.GasLimit,
			Value:    auth.Value,
		},
	}

	return session, nil
}

func getReceipt(tx *types.Transaction, useTimeout bool) (receipt *types.Receipt, err error) {
	ctxB := context.Background()
	var ctx context.Context
	var cancel context.CancelFunc

	if useTimeout {
		ctx, cancel = context.WithTimeout(ctxB, 60 * time.Minute)
	} else {
		ctx, cancel = context.WithCancel(ctxB)
	}

	defer cancel()

	receipt, err = bind.WaitMined(ctx, client, tx)
	if ctx.Err() != nil && err == ctx.Err() {
		return nil, errors.New(uErr.ErrTXTimedOut)
	} else if err != nil {
		return nil, err
	}

	return
}
