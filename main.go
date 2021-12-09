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

type Crypto []string

func (e Crypto) Find(symbol string) bool {
	for i := 0; i < len(e); i++ {
		if e[i] == symbol {
			return true
		}
	}
	return false
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

	wishlist := Crypto{"ADA", "BNB", "CRV", "DOT", "ETH", "EVX", "KUB", "POW", "WAN", "XLM", "XRP"}

	market, err := bk.MarketBalances()
	if err != nil {
		panic(err)
	}

	symbols, err := bk.MarketSymbols()
	if err != nil {
		panic(err)
	}

	var balanceTotal float64 = market["THB"].Available
	appLog.Printf(" %s - %.4f Baht", " THB", market["THB"].Available)
	for _, v := range symbols {
		coins := strings.Split(v.Symbol, "_")

		if !wishlist.Find(coins[1]) {
			continue
		}
		bl := market[coins[1]]

		ticker, err := bk.MarketTicker(v.Symbol)
		if err != nil {
			panic(err)
		}

		crypto := ticker[v.Symbol]
		balanceTotal += bl.Available * crypto.LowestAsk

		appLog.Printf(" %s - %.4f Baht", coins[1], bl.Available*crypto.LowestAsk)
	}
	appLog.Printf("Total Balacne: %.4f", balanceTotal)
}
