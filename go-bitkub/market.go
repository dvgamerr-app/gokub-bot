package bitkub

import (
	"encoding/json"
	"fmt"
)

type Balance struct {
	Available float64 `json:"available"`
	Reserved  float64 `json:"reserved"`
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
