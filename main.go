package main

import (
	"fmt"
	"gokub/go-bitkub"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"

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

	appLog = log.New(os.Stdout, " [Debug] ", log.Ltime)
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

func round(x float64) float64 {
	return math.Round((x)*1000) / 1000
}

func main() {
	appLog.Printf("Starting... (%s)", appTitle)
	bk := &bitkub.Config{ApiKey: os.Getenv(_APIKEY), SecretKey: os.Getenv(_SECRETKEY)}

	appLog.Println("Check status server...")
	if err := bk.GetStatus(); err != nil {
		panic(fmt.Sprintf(" - API::%s", err.Error()))
	}

	// serverTime, err := bk.GetServerTime()
	// if err != nil {
	// 	panic(fmt.Sprintf(" - API::%s", err.Error()))
	// }
	// appLog.Println("- Server Time:", serverTime.Format(time.RFC1123Z))
	// appLog.Println("-  Local Time:", time.Now().Format(time.RFC1123Z))

	wishlist := Crypto{"ADA", "BNB", "CRV", "DOT", "ETH", "EVX", "KUB", "POW", "WAN", "XLM", "XRP"}
	// wishlist := Crypto{"BTC", "KUB"}

	market, err := bk.Balances()
	if err != nil {
		panic(err)
	}

	symbols, err := bk.Symbols()
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

		ticker, err := bk.Ticker(v.Symbol)
		if err != nil {
			panic(err)
		}

		crypto := ticker[v.Symbol]
		balanceTotal += bl.Available * crypto.LowestAsk

		appLog.Printf(" %s - %.4f Baht", coins[1], bl.Available*crypto.LowestAsk)

		limit := 100
		page := 1
		history, err := bk.MyOrderHistory(v.Symbol, &page, &limit, nil, nil)
		if err != nil {
			panic(err)
		}

		totalFee := 0.0
		totalProfit := 0.0
		totalWallet := 0.0
		for _, o := range history {
			if o.Side == "sell" {
				totalProfit += round(o.Amount*o.Rate) - o.Fee
				totalWallet -= o.Amount

				// appLog.Printf("%v +%.2f (fee:%.2f)", o.Date.Format(time.RFC3339), round(o.Amount*o.Rate)-o.Fee, o.Fee)
			} else if o.Side == "buy" {
				totalProfit -= round(o.Amount*o.Rate) + o.Fee
				totalWallet += o.Amount

				// appLog.Printf("%v -%.2f (fee:%.2f)", o.Date.Format(time.RFC3339), round(o.Amount*o.Rate)+o.Fee, o.Fee)
			}
			totalFee += o.Fee
		}

		appLog.Printf("Wallet: %.8f fee: %f Profit: %f", totalWallet, totalFee, totalProfit)
	}
	// appLog.Printf("Total Balacne: %.4f", balanceTotal)
}
