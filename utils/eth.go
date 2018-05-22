package utils

import "xox_tokens/constants"

func GetInfuraProviderUrl(networkId int, token string) (url string) {
	switch networkId {
	case constants.MAINNET:
		url = "https://mainnet.infura.io/"
	case constants.ROPSTEN:
		url = "https://ropsten.infura.io/"
	case constants.RINKEBY:
		url = "https://rinkeby.infura.io/"
	case constants.KOVAN:
		url = "https://kovan.infura.io/"
	case constants.INFURANET:
		url = "https://infuranet.infura.io/"
	default:
		url = "https://rinkeby.infura.io/"
	}

	url += token
	return
}
