package websockets

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Writer commented
func ServerTime(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()

	for {
		serverTimeBytes := time.Now().Format(time.UnixDate)
		if err := ws.WriteMessage(websocket.TextMessage, []byte(serverTimeBytes)); err != nil {
			log.Printf("txtMsg: " + err.Error())
			return
		}
		time.Sleep(1 * time.Second)
	}

}
