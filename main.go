package main

import (
	"fmt"
	"gokub/go-bitkub"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tmilewski/goenv"
)

const (
	_ENV       = "ENV"
	_VERSION   = "VERSION"
	_APIKEY    = "BITKUB_API"
	_SECRETKEY = "BITKUB_SECRET"
)

var (
	appName    string = "gokub-bot"
	appVersion string = ""
	appTitle   string = ""
	appLog     *log.Logger
)

func init() {
	appIsProduction := os.Getenv(_ENV) == "production"
	if !appIsProduction {
		goenv.Load()
	}

	if os.Getenv(_APIKEY) == "" {
		panic(fmt.Sprintf("Environment name '%s' is empty.", _APIKEY))
	}
	if os.Getenv(_SECRETKEY) == "" {
		panic(fmt.Sprintf("Environment name '%s' is empty.", _SECRETKEY))
	}

	content, err := ioutil.ReadFile(_VERSION)
	if err != nil {
		content, _ = ioutil.ReadFile(fmt.Sprintf("../%s", _VERSION))
	}
	appVersion = strings.TrimSpace(string(content))
	appTitle = fmt.Sprintf("%s@%s", appName, appVersion)

	appLog = log.New(os.Stdout, "", log.Lshortfile|log.Ltime)
}
func main() {
	appLog.Printf("Starting... (%s)", appTitle)
	bk := &bitkub.Config{ApiKey: os.Getenv(_APIKEY), SecretKey: os.Getenv(_SECRETKEY)}

	appLog.Println("Check status server...")
	if err := bk.GetStatus(); err != nil {
		panic(fmt.Sprintf(" - API::%s", err.Error()))
	}

	serverTime, err := bk.GetServerTime()
	if err != nil {
		panic(fmt.Sprintf(" - API::%s", err.Error()))
	}
	appLog.Println("- Server Time:", serverTime.Format(time.RFC1123Z))

	market, err := bk.MarketBalances()
	if err != nil {
		panic(err)
	}

	for currency, v := range market {
		if v.Available == 0 && v.Reserved == 0 {
			continue
		}
		appLog.Printf("%s  = %f (%f)", currency, v.Available, v.Reserved)
	}

	symbols, err := bk.MarketSymbols()
	if err != nil {
		panic(err)
	}
	inSymbol := []string{"THB_BCH", "THB_KUB", "THB_OMG", "THB_BTC"}
	for _, v := range symbols {
		if len(inSymbol) > 0 {
			var include bool
			for _, in := range inSymbol {
				if v.Symbol == in {
					include = true
					break
				}
			}
			if !include {
				continue
			}
		}

		ticker, err := bk.MarketTicker(v.Symbol)
		if err != nil {
			panic(err)
		}

		appLog.Printf("%s|%d \tBid: %f Ask: %f", v.Symbol, v.ID, ticker[v.Symbol].HighestBid, ticker[v.Symbol].LowestAsk)
	}
}
