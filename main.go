package main

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

	"github.com/tmilewski/goenv"
)

func errorMessage(code int) string {
	switch code {
	case 1:
		return "Invalid JSON payload"
	case 2:
		return "Missing X-BTK-APIKEY"
	case 3:
		return "Invalid API key"
	case 4:
		return "API pending for activation"
	case 5:
		return "IP not allowed"
	case 6:
		return "Missing / invalid signature"
	case 7:
		return "Missing timestamp"
	case 8:
		return "Invalid timestamp"
	case 9:
		return "Invalid user"
	case 10:
		return "Invalid parameter"
	case 11:
		return "Invalid symbol"
	case 12:
		return "Invalid amount"
	case 13:
		return "Invalid rate"
	case 14:
		return "Improper rate"
	case 15:
		return "Amount too low"
	case 16:
		return "Failed to get balance"
	case 17:
		return "Wallet is empty"
	case 18:
		return "Insufficient balance"
	case 19:
		return "Failed to insert order into db"
	case 20:
		return "Failed to deduct balance"
	case 21:
		return "Invalid order for cancellation"
	case 22:
		return "Invalid side"
	case 23:
		return "Failed to update order status"
	case 24:
		return "Invalid order for lookup (or cancelled)"
	case 25:
		return "KYC level 1 is required to proceed"
	case 30:
		return "Limit exceeds"
	case 40:
		return "Pending withdrawal exists"
	case 41:
		return "Invalid currency for withdrawal"
	case 42:
		return "Address is not in whitelist"
	case 43:
		return "Failed to deduct crypto"
	case 44:
		return "Failed to create withdrawal record"
	case 45:
		return "Nonce has to be numeric"
	case 46:
		return "Invalid nonce"
	case 47:
		return "Withdrawal limit exceeds"
	case 48:
		return "Invalid bank account"
	case 49:
		return "Bank limit exceeds"
	case 50:
		return "Pending withdrawal exists"
	case 51:
		return "Withdrawal is under maintenance"
	case 90:
		return "Server error (please contact support)"
	default:
		return ""
	}
}

const (
	_ENV     = "ENV"
	_VERSION = "VERSION"
)

var (
	appName    string = "gokub-bot"
	appVersion string = ""
	appTitle   string = ""
)

const (
	APIKEY    = "BITKUB_API"
	SECRETKEY = "BITKUB_SECRET"
)

func init() {
	appIsProduction := os.Getenv(_ENV) == "production"
	if !appIsProduction {
		goenv.Load()
	}
	content, err := ioutil.ReadFile(_VERSION)
	if err != nil {
		content, _ = ioutil.ReadFile(fmt.Sprintf("../%s", _VERSION))
	}
	appVersion = strings.TrimSpace(string(content))
	appTitle = fmt.Sprintf("%s@%s", appName, appVersion)
}

func main() {

	// url := "https://notice.touno.io/line/popcorn/kem"
	// method := "PUT"

	ts := time.Now().UnixMilli()

	// hash := sha256.Sum256([]byte(fmt.Sprintf(`{"ts": %d}`, ts)))
	// sig := fmt.Sprintf("%x", hash[:])

	// mac := hmac.New(sha1.New, []byte(secretToken))
	// mac.Write([]byte(payloadBody))
	// expectedMAC := mac.Sum(nil)
	// return "sha1=" + hex.EncodeToString(expectedMAC)
	mac := hmac.New(sha256.New, []byte(os.Getenv(SECRETKEY)))
	mac.Write([]byte(fmt.Sprintf(`{"ts":%d}`, ts)))
	sig := hex.EncodeToString(mac.Sum(nil))

	payload := strings.NewReader(fmt.Sprintf(`{"ts":%d,"sig":"%s"}`, ts, sig))

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.bitkub.com/api/market/wallet", payload)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-BTK-APIKEY", os.Getenv(APIKEY))

	// fmt.Println("payload", payload)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	errorCode := int(dat["error"].(float64))

	if errorMessage(errorCode) != "" {
		panic(fmt.Sprintf("%d - %s", errorCode, errorMessage(errorCode)))
	}

	fmt.Println(dat)
}
