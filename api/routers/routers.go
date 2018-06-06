package routers

import (
	"math/big"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"log"
	"fmt"

	uErr "xox_tokens/errors"
	"xox_tokens/crypto/eth"
	"xox_tokens/db/models/buyer"
	"xox_tokens/crypto/btc"
	"xox_tokens/configs"
	"xox_tokens/bus"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/charge"
	"github.com/ethereum/go-ethereum/common"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func RoutersStart(conf *configs.Configs, envType string) {
	log.Println("Routers starting")
	stripe.Key = conf.Stripe[envType].SecretKey
	router := gin.Default()
	router.HTMLRender = createMyRender()

	router.Use(cors.Default())

	router.GET("/api/v1/stripe/", func(c *gin.Context) {
		value, isValue := c.GetQuery("value")
		ethAddress, isAddress := c.GetQuery("ethAddress")
		valint, _ := strconv.Atoi(value)
		isError := false
		errMessage := ""
		if !isValue {
			errMessage += "Is no value. "
			isError = true
		} else if valint <= 0 {
			errMessage += "Not valid value. "
			isError = true
		}
		if !isAddress {
			errMessage += "Is no address. "
			isError = true
		} else if !common.IsHexAddress(ethAddress) {
			errMessage += "Not valid address. "
			isError = true
		}
		if isError {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error": errMessage,
			})
			return
		}
		title := c.DefaultQuery("title", "XOX tokens")
		description := c.DefaultQuery("desc", "Payment description")
		c.HTML(http.StatusOK, "form", gin.H{
			"title": title,
			"key": conf.Stripe[envType].PublishableKey,
			"amount": value,
			"ethAddress": ethAddress,
			"description": description,
		})

	})

	router.POST("/api/v1/stripe/request", func(c *gin.Context) {
		type Resp struct {
			Token string
			Email string
			EthAddress string
			Amount uint64
		}
		var StripeResp map[string]interface{}
		var resp Resp
		err := c.BindJSON(&resp)

		if err != nil {
			uErr.Fatal(err, "JSON parse error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error": "JSON parse error",
			})
			return
		}

		isError := false
		errMessage := ""
		if resp.Amount <= 0 {
			errMessage += "Not valid value. "
			isError = true
		}
		if !common.IsHexAddress(resp.EthAddress) {
			errMessage += "Not valid address. "
			isError = true
		}
		if isError {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error": errMessage,
			})
			return
		}

		customerParams := &stripe.CustomerParams{Email: resp.Email}
		customerParams.SetSource(resp.Token)

		newCustomer, err := customer.New(customerParams)

		if err != nil {
			uErr.Fatal(err, http.StatusInternalServerError)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error": err, // TODO: Edit this?
			})
			return
		}

		chargeParams := &stripe.ChargeParams{
			Amount: resp.Amount,
			Currency: "usd",
			Desc:     "Sample Charge",
			Customer: newCustomer.ID,
		}

		ch, err := charge.New(chargeParams)
		if err != nil {
			json.Unmarshal([]byte(err.Error()), &StripeResp)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": StripeResp["status"],
				"error": StripeResp["message"],
			})
			return
		}

		rate, _ := eth.GetRate()
		ethVal := big.NewInt(int64(rate.Uint64() * ch.Amount) / 10)
		eth.GetValueForBackend(ethVal, common.HexToAddress(resp.EthAddress))

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"customer": newCustomer.ID,
			"amount": ch.Amount,
			"ethAddress": resp.EthAddress,
		})

	})

	router.POST("/api/v1/btc/GenerateKey", func(c *gin.Context) {
		type Resp struct {
			EthAddress string
			BtcAddress string
		}
		var resp Resp

		if err := c.BindJSON(&resp); err != nil {
			uErr.Fatal(err, "JSON parse error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error": "JSON parse error",
			})
			return
		}

		// Validate ETH address.
		isValidEthAddress := true
		if resp.EthAddress == "" {
			isValidEthAddress = false
		}
		if !common.IsHexAddress(resp.EthAddress) {
			isValidEthAddress = false
		}
		if common.EmptyHash(common.HexToHash(resp.EthAddress)) {
			isValidEthAddress = false
		}
		ethAddr := common.HexToAddress(resp.EthAddress)
		if strings.ToLower(resp.EthAddress) != strings.ToLower(ethAddr.Hex()) {
			isValidEthAddress = false
		}
		if !isValidEthAddress {
			log.Println(uErr.ErrValidateETHAddress)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error": uErr.ErrValidateETHAddress,
			})
			return
		}


		// Define net configs.
		chainConf := &chaincfg.MainNetParams
		isTestnet := configs.GetConfigs().Crypto.BTCTestnet
		if isTestnet {
			chainConf = &chaincfg.TestNet3Params
		}

		// Validate BTC address.
		isValidBtcAddress := true
		if resp.BtcAddress == "" {
			isValidBtcAddress = false
		}
		_, err := btcutil.DecodeAddress(resp.BtcAddress, chainConf)
		if err != nil {
			isValidBtcAddress = false
		}
		if !isValidBtcAddress {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error": uErr.ErrValidateBTCAddress,
			})
			return
		}


		// Check if buyer exist.
		buyer_, found, err := buyer.FindByETHAddress(resp.EthAddress)
		if err != nil {
			log.Println(uErr.ErrSelectFromDB)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error": uErr.ErrSelectFromDB,
			})
			return
		}

		var wif, addr string

		// Create new buyer if doesn't exist.
		if !found {
			wif, addr, err := btc.GenerateKey()
			if err != nil {
				log.Println(uErr.ErrGenerateAddress)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status": http.StatusBadRequest,
					"error": uErr.ErrGenerateAddress,
				})
				return
			}
			buyer_, err = buyer.New(resp.EthAddress, addr, wif)
			if err != nil {
				log.Println(uErr.ErrInsertToDB)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status": http.StatusBadRequest,
					"error": uErr.ErrInsertToDB,
				})
				return
			}
		}

		wif = buyer_.BtcAddrWif
		addr = buyer_.BtcAddr
		log.Println(wif)

		bc, err := bus.Dial(fmt.Sprintf("http://0.0.0.0:%v", conf.Server.NATSPort))
		defer bc.Close()

		bc.Publish(bus.SubjectAddress, &bus.EventAddress{Address: addr})

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			//"wif": wif,
			"addr": addr,
		})
	})

	router.Run(":8080")
	log.Println("Routers starts")
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("form", "templates/base.tmpl", "templates/form.tmpl")
	r.AddFromFiles("body", "templates/base.tmpl", "templates/body.tmpl")
	r.AddFromFiles("article", "templates/base.html", "templates/index.html", "templates/article.html")
	return r
}
