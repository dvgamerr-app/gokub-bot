package bitkub

import (
	"encoding/json"
	"fmt"
	"strconv"
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

	i, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		panic(err)
	}

	return time.Unix(i, 0), nil
}
