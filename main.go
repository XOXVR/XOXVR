// Need to update .bashrc
// export MICRO_REGISTRY=mdns

package main

import (
	"flag"
	"os"
	"log"
	"fmt"
	"context"
	"time"
	"strconv"
	"math/big"
	"net/http"
	"encoding/json"

	uErr "xox_tokens/errors"
	"xox_tokens/db"
	"xox_tokens/utils"
	"xox_tokens/configs"
_	"xox_tokens/crypto/btc"
	"xox_tokens/crypto/eth"

	"github.com/micro/go-micro"
	"xox_tokens/services/ExchangeRateService/proto"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/multitemplate"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/gin-contrib/cors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/micro/cli"
)

func main() {
	prod := flag.Bool("prod", false, "Run in production mode.")
	flag.Parse()

	// Run watchers.
	defer utils.RecoverWatcher(shutdown)
	go utils.ShutdownWatcher(shutdown)

	// Parse configs.
	conf, err := configs.ParseConfigs("./configs.toml")
	if err != nil {
		uErr.Fatal(err, "failed to parse configs")
	}
	conf, err = configs.ParseConfigs("./configs_local.toml")
	if err != nil {
		uErr.Fatal(err, "failed to parse local configs")
	}

	// Setup prod env.
	envType := "dev"
	if *prod {
		envType = "prod"
		conf.Common.Dev = false

		err := utils.SetupLogFile(conf.Common.LogOutPath)
		if err != nil {
			uErr.Fatal(err, "failed to setup log file")
		}

		gin.SetMode(gin.ReleaseMode)
	}

	fmt.Println("==> " + envType) // TODO: Delete this line

	// Initiate database.
	err = db.Initiate(conf.DB[envType])
	if err != nil {
		uErr.Fatal(err, "failed to initialize db")
	}
	defer db.Instance.Close()

	// Start connection with ETH provider.
	err = eth.Dial(&conf.Crypto)
	if err != nil {
		uErr.Fatal(err, "failed ETH Dial")
	}

	go eth.Refund()

	go func() {
		// Create a new service. Optionally include some options here.
		service := micro.NewService(
			micro.Name("Currency.client"),
			micro.Flags(
				cli.BoolFlag{
					Name:  "prod",
					Usage: "This is a prod flag",
				},
			),
		)
		service.Init(
		)

		// Create new greeter client
		greeter := proto.NewCurrencyService("Currency", service.Client())

		for {
			// Call the greeter
			rsp, err := greeter.GetRate(context.TODO(), &proto.GetRateReq{BaseCurrency: "FiatETH"})
			if err != nil {
				log.Println(err)
			}

			rate, err := strconv.ParseFloat(rsp.CurrencyRatio, 64)
			if err != nil {
				log.Println(err)
			}

			irate := int64(rate * float64(1e18))

			eth.SetRate(big.NewInt(irate))

			time.Sleep(time.Hour * 12)
		}

	}()

	// Template for REST
	go func() {
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
			if isPostICO, err := eth.IsPostICOStage(); err != nil {
				uErr.Fatal(err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status": http.StatusBadRequest,
					"error": err.Error(),
				})
				return
			} else if !isPostICO {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status": http.StatusBadRequest,
					"error": "Not Post-ICO Stage",
				})
				return
			}

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

		router.Run(":8080")
	}()

	select {}
}

func shutdown(fatal bool, r interface{}) {
	if fatal {
		os.Exit(1)
	}

	log.Println("Shutdown")
	os.Exit(0)
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("form", "templates/base.tmpl", "templates/form.tmpl")
	r.AddFromFiles("body", "templates/base.tmpl", "templates/body.tmpl")
	r.AddFromFiles("article", "templates/base.html", "templates/index.html", "templates/article.html")
	return r
}
