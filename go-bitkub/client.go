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

func getHashSignature(secretKey string) (int64, string) {
	ts := time.Now().UnixMilli()
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(fmt.Sprintf(`{"ts":%d}`, ts)))
	return ts, hex.EncodeToString(mac.Sum(nil))
}

func createClientHTTP(cfg *Config, method string, path string, secure bool) ([]byte, error) {
	var payload *strings.Reader = strings.NewReader("{}")

	if secure {
		ts, sig := getHashSignature(cfg.SecretKey)
		payload = strings.NewReader(fmt.Sprintf(`{"ts":%d,"sig":"%s"}`, ts, sig))
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
