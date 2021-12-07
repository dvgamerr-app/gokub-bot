package bitkub

import (
	"encoding/json"
	"fmt"
)

type Symbols struct {
	ID     int64  `json:"id"`
	Symbol string `json:"symbol"`
	Info   string `json:"info"`
}

func (cfg *Config) MarketSymbols() ([]*Symbols, error) {
	body, err := createClientHTTP(cfg, "GET", _API_MARKET_SYMBOLS, true)

	if err != nil {
		return nil, fmt.Errorf("'%s' %+v", _API_MARKET_SYMBOLS, err)
	}

	var res ResponseArray
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Body:", string(body))
		return nil, fmt.Errorf("'%s' %+v", _API_MARKET_SYMBOLS, err)
	}

	result := []*Symbols{}

	for _, val := range res.Result {
		fields, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("MarketSymbols:: %+v", val)
		}

		id, ok := fields["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("MarketSymbols:: %+v", val)
		}

		result = append(result, &Symbols{
			ID:     int64(id),
			Symbol: fields["symbol"].(string),
			Info:   fields["info"].(string),
		})
	}

	return result, nil
}

type Balance struct {
	Available float64 `json:"available"`
	Reserved  float64 `json:"reserved"`
}

func (cfg *Config) MarketBalances() (map[string]*Balance, error) {
	body, err := createClientHTTP(cfg, "POST", _API_MARKET_BALANCES, true)
	if err != nil {
		return nil, fmt.Errorf("'%s' %+v", _API_MARKET_BALANCES, err)
	}

	var res ResponseKeyValues
	if err := res.Unmarshal(body); err != nil {
		return nil, fmt.Errorf("'%s' %+v", _API_MARKET_BALANCES, err)
	}

	result := map[string]*Balance{}

	for currency, val := range res.Result {
		fields, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("MarketBalances:: %+v", val)
		}

		result[currency] = new(Balance)
		result[currency].Available = fields["available"].(float64)
		result[currency].Reserved = fields["reserved"].(float64)
	}
	return result, nil
}
