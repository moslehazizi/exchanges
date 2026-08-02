package bitget

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"order-maker/internal/bgsign"
	"order-maker/internal/env"
)

const (
	baseURL     = "https://api.bitget.com"
	requestPath = "/api/v3/trade/cancel-order"
	requestBody = `{"orderId":"1467791656000839680", "clientOid":"1234567890", "category":"USDT-FUTURES"}`
)

func Run() {
	env.Load()

	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := bgsign.Sign(secretKey, timestamp, "POST", requestPath, requestBody)

	req, err := http.NewRequest(http.MethodPost, baseURL+requestPath, bytes.NewReader([]byte(requestBody)))
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
