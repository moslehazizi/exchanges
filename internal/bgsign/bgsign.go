package bgsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// Sign produces the Bitget ACCESS-SIGN header value. requestPath must
// already include the query string (e.g. "/api/v3/p2p/balance?token=USDT")
// when the request has one.
func Sign(secretKey, timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
