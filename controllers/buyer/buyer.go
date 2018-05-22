package buyer

import (
	"errors"
	"strings"
	"log"

	uErr "xox_tokens/errors"
	buyerGW "xox_tokens/api/buyer"
	"xox_tokens/configs"
	"xox_tokens/db/models/buyer"

	"golang.org/x/net/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

type RPC struct {}

// AddBuyer adds new buyer.
func (r *RPC) GetBTCWallet(ctx context.Context, in *buyerGW.GetBTCWalletReq) (*buyerGW.GetBTCWalletResp, error) {
	ethAddrStr := in.EthAddress
	btcAddrStr := in.BtcAddress

	// Validate ETH address.
	if ethAddrStr == "" {
		return nil, errors.New(uErr.ErrValidateETHAddress)
	}
	if !common.IsHexAddress(ethAddrStr) {
		return nil, errors.New(uErr.ErrValidateETHAddress)
	}
	if common.EmptyHash(common.HexToHash(ethAddrStr)) {
		return nil, errors.New(uErr.ErrValidateETHAddress)
	}
	ethAddr := common.HexToAddress(ethAddrStr)
	if strings.ToLower(ethAddrStr) != strings.ToLower(ethAddr.Hex()) {
		return nil, errors.New(uErr.ErrValidateETHAddress)
	}

	// Define net configs.
	chainConf := &chaincfg.MainNetParams
	isTestnet := configs.GetConfigs().Crypto.BTCTestnet
	if isTestnet {
		chainConf = &chaincfg.TestNet3Params
	}

	// Validate BTC address.
	if btcAddrStr == "" {
		return nil, errors.New(uErr.ErrValidateBTCAddress)
	}
	_, err := btcutil.DecodeAddress(btcAddrStr, chainConf)
	if err != nil {
		return nil, errors.New(uErr.ErrValidateBTCAddress)
	}

	// Check if buyer exist.
	_, found, err := buyer.FindByBTCAddress(btcAddrStr)
	if err != nil {
		log.Println(uErr.ErrSelectFromDB)
		return nil, uErr.CombineErrCode(uErr.ErrSelectFromDBCode)
	}

	// Create new buyer if doesn't exist.
	if !found {
		err = buyer.New(ethAddrStr, btcAddrStr)
		if err != nil {
			log.Println(uErr.ErrInsertToDB)
			return nil, uErr.CombineErrCode(uErr.ErrInsertToDBCode)
		}
	}

	walletAddr := configs.GetConfigs().Crypto.BTCAddr
	return &buyerGW.GetBTCWalletResp{BtcAddress: walletAddr}, nil
}
