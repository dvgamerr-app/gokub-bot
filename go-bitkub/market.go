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

type Ticker struct {
	ID            int64   `json:"id"`
	Last          float64 `json:"last"`
	LowestAsk     float64 `json:"lowestAsk"`
	HighestBid    float64 `json:"highestBid"`
	PercentChange float64 `json:"percentChange"`
	BaseVolume    float64 `json:"baseVolume"`
	QuoteVolume   float64 `json:"quoteVolume"`
	IsFrozen      bool    `json:"isFrozen"`
	High24Hr      float64 `json:"high24hr"`
	Low24Hr       float64 `json:"low24hr"`
	Change        float64 `json:"change"`
	PrevClose     float64 `json:"prevClose"`
	PrevOpen      float64 `json:"prevOpen"`
}

func (cfg *Config) MarketTicker(symbol ...string) (map[string]*Ticker, error) {
	url := _API_MARKET_TICKER
	if len(symbol) > 0 {
		url = fmt.Sprintf("%s?sym=%s", _API_MARKET_TICKER, symbol[0])
	}
	body, err := createClientHTTP(cfg, "GET", url, true)

	if err != nil {
		return nil, fmt.Errorf("'%s' %+v", url, err)
	}

	var res map[string]map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Body:", string(body))
		return nil, fmt.Errorf("'%s' %+v", url, err)
	}

	result := map[string]*Ticker{}

	for sb, val := range res {
		id, ok := val["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("MarketTicker:: %+v", val)
		}

		frozen, ok := val["isFrozen"].(float64)
		if !ok {
			return nil, fmt.Errorf("MarketTicker:: %+v", val)
		}

		result[sb] = &Ticker{
			ID:            int64(id),
			Last:          val["last"].(float64),
			LowestAsk:     val["lowestAsk"].(float64),
			HighestBid:    val["highestBid"].(float64),
			PercentChange: val["percentChange"].(float64),
			BaseVolume:    val["baseVolume"].(float64),
			QuoteVolume:   val["quoteVolume"].(float64),
			IsFrozen:      frozen == 1,
			High24Hr:      val["high24hr"].(float64),
			Low24Hr:       val["low24hr"].(float64),
			Change:        val["change"].(float64),
			PrevClose:     val["prevClose"].(float64),
			PrevOpen:      val["prevOpen"].(float64),
		}
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

		result[currency] = &Balance{
			Available: fields["available"].(float64),
			Reserved:  fields["reserved"].(float64),
		}
	}
	return result, nil
}
