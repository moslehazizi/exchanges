package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL     = "https://api.bitget.com"
	requestPath = "/api/v3/account/settings"
)

func loadDotEnv() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		path := filepath.Join(dir, ".env")
		if data, err := os.ReadFile(path); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) != 2 {
					continue
				}
				os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

func sign(timestamp, method, requestPath, body, secretKey string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	loadDotEnv()

	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := sign(timestamp, "GET", requestPath, "", secretKey)

	url := baseURL + requestPath

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
