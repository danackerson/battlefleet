package routes

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/gorilla/websocket"
)

// ServeWebSocket serves web socket connections
func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	serverPort := app.HTTPPort
	remoteHostSettings := strings.Split(r.Host, ":")
	if len(remoteHostSettings) > 1 {
		serverPort = remoteHostSettings[1]
	}

	// test server runs on different port
	if serverPort == app.HTTPPort && r.Header.Get("Origin") != app.URIScheme+"://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true, // NOTE: Traefik has bug that doesn't pass thru :(
	}

	ws, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println("error: " + err.Error())
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println("WS handshake + " + err.Error())
		}
		return
	}

	go retrieveGame(ws, r)
}

func retrieveGame(ws *websocket.Conn, r *http.Request) {
	defer ws.Close()

	for {
		// even though we don't do much with the Reads (yet)
		// it's helpful to negotiate client side CLOSE requests
		// msg types: https://github.com/gorilla/websocket/blob/master/conn.go#L61
		messageType, msg, err := ws.ReadMessage()
		log.Printf("%d: %s (ERR: %v)", messageType, msg, err)

		if messageType != -1 && messageType != 8 {
			session, sessionErr := app.SessionStore.Get(r, app.SessionCookieKey)
			if sessionErr != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("ERR:"+sessionErr.Error()))
				return
			}
			account := getAccount(r, session)
			if account != nil {
				game := account.GetGame()
				game.LastTurn = time.Now()
				game.Map = nil // save ~400KB of useless crap on the client side PER WS request
				ws.WriteJSON(game)
			} else {
				ws.WriteMessage(websocket.TextMessage, []byte(time.Now().Format(time.UnixDate)))
			}
		} else {
			log.Printf("Client hung up...good-bye!\n")
			return
		}

		time.Sleep(1 * time.Second)
	}
}
