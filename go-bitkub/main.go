package bitkub

import (
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	body, err := createClientHTTP(cfg, "GET", _API_STATUS, nil)
	if err != nil {
		return fmt.Errorf("'%s' %+v", _API_STATUS, err)
	}

	var data ResponseKeyValues
	if err := data.Unmarshal(body); err == nil {
		return fmt.Errorf("'%s' %+v", _API_STATUS, data.GetError())
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
	body, err := createClientHTTP(cfg, "GET", _API_SERVERTIME, nil)

	if err != nil {
		return time.Time{}, fmt.Errorf("'%s' %+v", _API_SERVERTIME, err)
	}

	i, err := strconv.ParseInt(string(body), 0, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("'%s' %+v", _API_SERVERTIME, err)
	}

	return time.Unix(i, 0), nil
}
