package bitkub

import (
	"encoding/json"
	"fmt"
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

func (cfg *Config) GetStatus() (map[string]*Balance, error) {
	body, err := newClientHTTP(cfg, "POST", _API_MARKET_BALANCES, false)

	if err != nil {
		return nil, err
	}

	var data APIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.IsError() {
		return nil, data.GetError()
	}

	result := map[string]*Balance{}

	for currency, val := range data.Result {
		fields, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("MarketBalances:: %+v", fields)
		}

		result[currency] = new(Balance)
		result[currency].Available = fields["available"].(float64)
		result[currency].Reserved = fields["reserved"].(float64)
	}
	return result, nil
}

func (cfg *Config) GetServerTime() (time.Time, error) {
	body, err := newClientHTTP(cfg, "POST", _API_MARKET_BALANCES, false)

	if err != nil {
		return time.Time{}, err
	}

	var data APIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return time.Time{}, err
	}

	if data.IsError() {
		return time.Time{}, data.GetError()
	}

	result := map[string]*Balance{}

	for currency, val := range data.Result {
		fields, ok := val.(map[string]interface{})
		if !ok {
			return time.Time{}, fmt.Errorf("MarketBalances:: %+v", fields)
		}

		result[currency] = new(Balance)
		result[currency].Available = fields["available"].(float64)
		result[currency].Reserved = fields["reserved"].(float64)
	}
	return time.Now(), nil
}
