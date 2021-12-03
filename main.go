package main

import (
	"fmt"
	"gokub/bitkub-api"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tmilewski/goenv"
)

const (
	_ENV     = "ENV"
	_VERSION = "VERSION"
)

var (
	appName    string = "gokub-bot"
	appVersion string = ""
	appTitle   string = ""
)

func init() {
	appIsProduction := os.Getenv(_ENV) == "production"
	if !appIsProduction {
		goenv.Load()
	}
	content, err := ioutil.ReadFile(_VERSION)
	if err != nil {
		content, _ = ioutil.ReadFile(fmt.Sprintf("../%s", _VERSION))
	}
	appVersion = strings.TrimSpace(string(content))
	appTitle = fmt.Sprintf("%s@%s", appName, appVersion)
}
func main() {
	fmt.Printf("%s starting...\n", appTitle)
	err := bitkub.MarketBalances()
	if err != nil {
		panic(err)
	}
}
