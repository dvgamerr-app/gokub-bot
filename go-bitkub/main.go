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
)

var (
	UrlAPI = "https://api.bitkub.com"
)

type Config struct {
	ApiKey    string
	SecretKey string
}

func (cfg *Config) Init() {
	if cfg.ApiKey == "" {
		panic("'ApiKey' is empty.")
	}
	if cfg.SecretKey == "" {
		panic("'SecretKey' is empty.")
	}
}

func getHashSignature(secretKey string) (int64, string) {
	ts := time.Now().UnixMilli()
	mac := hmac.New(sha256.New, []byte(os.Getenv(secretKey)))
	mac.Write([]byte(fmt.Sprintf(`{"ts":%d}`, ts)))
	return ts, hex.EncodeToString(mac.Sum(nil))
}

func newClientHTTP(cfg *Config, method string, path string) ([]byte, error) {
	ts, sig := getHashSignature(cfg.SecretKey)
	payload := strings.NewReader(fmt.Sprintf(`{"ts":%d,"sig":"%s"}`, ts, sig))

	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", UrlAPI, path), payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-BTK-APIKEY", os.Getenv(cfg.ApiKey))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (cfg *Config) MarketBalances() (map[string]*Balance, error) {
	body, err := newClientHTTP(cfg, "POST", _API_MARKET_BALANCES)

	if err != nil {
		return nil, err
	}

	var data APIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.IsError() {
		return nil, fmt.Errorf("%d - %s", data.Error, errorMessage(data.Error))
	}

	result := map[string]*Balance{}

	for currency, val := range data.Result {
		fields, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%+v", fields)
		}

		result[currency] = new(Balance)
		result[currency].Available = fields["available"].(float64)
		result[currency].Reserved = fields["reserved"].(float64)
	}
	return result, nil
}

type Balance struct {
	Available float64 `json:"available"`
	Reserved  float64 `json:"reserved"`
}
