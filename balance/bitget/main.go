package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	apiKey     = ""
	secretKey  = ""
	passphrase = ""

	baseURL     = "https://api.bitget.com"
	requestPath = "/api/v3/p2p/balance"
	queryString = "token=USDT"
)

func sign(timestamp, method, requestPath, queryString, body string) string {
	message := timestamp + method + requestPath + "?" + queryString + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := sign(timestamp, "GET", requestPath, queryString, "")

	url := baseURL + requestPath + "?" + queryString

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("ACCESS-KEY", apiKey)
	req.Header.Set("ACCESS-SIGN", signature)
	req.Header.Set("ACCESS-PASSPHRASE", passphrase)
	req.Header.Set("ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("locale", "en-US")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("paptrading", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(string(respBody))
}
