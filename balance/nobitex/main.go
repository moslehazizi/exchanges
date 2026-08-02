package nobitex

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"order-maker/internal/env"
)

func decodeURLBase64(s string) ([]byte, error) {
	if b, err := base64.URLEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	return base64.RawURLEncoding.DecodeString(s)
}

func Run() {
	env.Load()

	publicKey := os.Getenv("NOBITEX_KEY")
	privateKeyB64 := os.Getenv("NOBITEX_PRIVATE_KEY")
	if publicKey == "" || privateKeyB64 == "" {
		fmt.Println("set NOBITEX_KEY and NOBITEX_PRIVATE_KEY environment variables")
		return
	}

	seed, err := decodeURLBase64(privateKeyB64)
	if err != nil {
		fmt.Println("invalid NOBITEX_PRIVATE_KEY:", err)
		return
	}
	privateKey := ed25519.NewKeyFromSeed(seed)

	method := "POST"
	fullPath := "/users/wallets/balance"
	body := `{"currency":"btc"}`

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signPayload := timestamp + method + fullPath + body
	signature := ed25519.Sign(privateKey, []byte(signPayload))
	signatureB64 := base64.URLEncoding.EncodeToString(signature)

	req, err := http.NewRequest(method, "https://apiv2.nobitex.ir"+fullPath, bytes.NewReader([]byte(body)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Nobitex-Key", publicKey)
	req.Header.Set("Nobitex-Signature", signatureB64)
	req.Header.Set("Nobitex-Timestamp", timestamp)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(respBody))
}
