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
	}

	ws, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println("error: " + err.Error())
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println("WS handshake + " + err.Error())
		}
		return
	}

	go serverTime(ws)
}

func serverTime(ws *websocket.Conn) {
	defer ws.Close()

	for {
		// even though we don't do much with the Reads (yet)
		// it's helpful to negotiate client side CLOSE requests
		// msg types: https://github.com/gorilla/websocket/blob/master/conn.go#L61
		messageType, msg, err := ws.ReadMessage()
		//log.Printf("%d: %s (ERR: %v)", messageType, msg, err)

		if messageType != -1 && messageType != 8 {
			serverTimeBytes := time.Now().Format(time.UnixDate)
			if err = ws.WriteMessage(websocket.TextMessage, []byte(serverTimeBytes)); err != nil {
				log.Printf("%d: %s (ERR: %s)", messageType, msg, err.Error())
				return
			}
		} else {
			log.Printf("Client hung up...good-bye!\n")
			return
		}

		time.Sleep(1 * time.Second)
	}
}
