package bitget

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"order-maker/internal/bgsign"
	"order-maker/internal/env"
)

const wsURL = "wss://wspap.bitget.com/v3/ws/private"

func Run() {
	env.Load()

	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// WS login uses a seconds timestamp, unlike the ms timestamp used by REST.
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := bgsign.Sign(secretKey, timestamp, "GET", "/user/verify", "")

	loginMsg := fmt.Sprintf(
		`{"op":"login","args":[{"apiKey":"%s","passphrase":"%s","timestamp":"%s","sign":"%s"}]}`,
		apiKey, passphrase, timestamp, signature,
	)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(loginMsg)); err != nil {
		panic(err)
	}

	_, loginResp, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(loginResp))

	subscribeMsg := `{"op":"subscribe","args":[{"instType":"UTA","topic":"account"}]}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg)); err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		fmt.Println(string(message))
	}
}
