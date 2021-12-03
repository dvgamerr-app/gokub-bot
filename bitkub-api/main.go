package bitkub

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tmilewski/goenv"
)

const (
	APIKEY    = "BITKUB_API"
	SECRETKEY = "BITKUB_SECRET"
	urlAPI    = "https://api.bitkub.com"
)

type APIResponse struct {
	Error  int                    `json:"error"`
	Result map[string]interface{} `json:"result"`
}

func init() {
	goenv.Load()
	if os.Getenv(APIKEY) == "" {
		panic(fmt.Sprintf("Environment name '%s' is empty.", APIKEY))
	}
	if os.Getenv(SECRETKEY) == "" {
		panic(fmt.Sprintf("Environment name '%s' is empty.", SECRETKEY))
	}
}

func getHashSignature() (int64, string) {
	ts := time.Now().UnixMilli()
	mac := hmac.New(sha256.New, []byte(os.Getenv(SECRETKEY)))
	mac.Write([]byte(fmt.Sprintf(`{"ts":%d}`, ts)))
	return ts, hex.EncodeToString(mac.Sum(nil))
}

func clientHTTP(method string, path string) error {
	ts, sig := getHashSignature()
	payload := strings.NewReader(fmt.Sprintf(`{"ts":%d,"sig":"%s"}`, ts, sig))

	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", urlAPI, path), payload)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-BTK-APIKEY", os.Getenv(APIKEY))

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data APIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if errorMessage(data.Error) != "" {
		return fmt.Errorf("%d - %s", data.Error, errorMessage(data.Error))
	}
	fmt.Println(data.Result)
	return nil
}

func MarketBalances() error {
	err := clientHTTP("POST", "/api/market/balances")
	if err != nil {
		return err
	}
	return nil
}
