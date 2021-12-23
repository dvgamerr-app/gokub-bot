package bitkub

import (
	"fmt"
)

type ResponseKeyValues struct {
	Error  int                    `json:"error"`
	Result map[string]interface{} `json:"result"`
}

func (e *ResponseKeyValues) IsError() bool {
	return e.Error != 0
}

func (e *ResponseKeyValues) GetError() error {
	return fmt.Errorf("%d - %s", e.Error, errorMessage(e.Error))
}

func (e *ResponseKeyValues) Unmarshal(body []byte) error {
	if err := json.Unmarshal(body, e); err != nil {
		return err
	}

	if e.IsError() {
		return e.GetError()
	}

	return nil
}

type ResponseArray struct {
	Error  int           `json:"error"`
	Result []interface{} `json:"result"`
}

func (e *ResponseArray) IsError() bool {
	return e.Error != 0
}

func (e *ResponseArray) GetError() error {
	return fmt.Errorf("%d - %s", e.Error, errorMessage(e.Error))
}

func (e *ResponseArray) Unmarshal(body []byte) error {
	if err := json.Unmarshal(body, e); err != nil {
		return err
	}

	if e.IsError() {
		return e.GetError()
	}

	return nil
}

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
