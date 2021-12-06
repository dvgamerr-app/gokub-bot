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

type StatusAPI struct {
	Secure    bool `json:"secure"`
	NonSecure bool `json:"non-secure"`
}

func (cfg *Config) GetStatus() error {
	body, err := newClientHTTP(cfg, "GET", _API_STATUS, false)
	if err != nil {
		return err
	}

	var data APIResponse
	if err := json.Unmarshal(body, &data); err == nil {
		return data.GetError(_API_STATUS)
	}

	var result []map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	for _, status := range result {
		if status["status"] != "ok" {
			return fmt.Errorf("%s - %s", status["name"], status["message"])
		}
	}

	return nil
}

func (cfg *Config) GetServerTime() (time.Time, error) {
	body, err := newClientHTTP(cfg, "GET", _API_SERVERTIME, false)

	if err != nil {
		return time.Time{}, err
	}

	var data APIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return time.Time{}, err
	}

	if data.IsError() {
		return time.Time{}, data.GetError(_API_SERVERTIME)
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
