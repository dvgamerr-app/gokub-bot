package bitkub

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type PayloadHMAC struct {
	Symbol *string `json:"sym,omitempty"`
	Amount *string `json:"amt,omitempty"`
	Rate   *string `json:"rat,omitempty"`
	Type   *string `json:"typ,omitempty"`

	ID      *string `json:"id,omitempty"`
	SellBuy *string `json:"sd,omitempty"`
	Hash    *string `json:"hash,omitempty"`

	Page  *int `json:"p,omitempty"`
	Limit *int `json:"lmt,omitempty"`

	Currency *int `json:"cur,omitempty"`
	Address  *int `json:"adr,omitempty"`
	Memo     *int `json:"mem,omitempty"`

	StartTimestamp *int64 `json:"start,omitempty"`
	EndTimestamp   *int64 `json:"end,omitempty"`

	ClientID *string `json:"client_id,omitempty"`

	Timestamp int64   `json:"ts"`
	Signature *string `json:"sig,omitempty"`
}

func (e *PayloadHMAC) getHashSignature(secretKey string) error {
	e.Timestamp = time.Now().UnixMilli()
	mac := hmac.New(sha256.New, []byte(secretKey))

	body, err := json.Marshal(e)
	if err != nil {
		return err
	}
	// fmt.Printf("   --- %+v ---\n", string(body))
	mac.Write(body)
	hash := hex.EncodeToString(mac.Sum(nil))
	e.Signature = &hash
	return nil
}

func createClientHTTP(cfg *Config, method string, path string, bodySign *PayloadHMAC) ([]byte, error) {
	payload := strings.NewReader("{}")

	if method == "POST" {
		// fmt.Printf("[%s] %s\n", method, fmt.Sprintf("%s%s", UrlAPI, path))

		err := bodySign.getHashSignature(cfg.SecretKey)
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(bodySign)
		if err != nil {
			return nil, err
		}

		// fmt.Printf("[%s] %+v\n", method, string(body))
		payload = strings.NewReader(string(body))
	} else {
		// fmt.Printf(" [%s] %s\n", method, fmt.Sprintf("%s%s", UrlAPI, path))
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", UrlAPI, path), payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-BTK-APIKEY", cfg.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
