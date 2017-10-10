package websockets

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// ServerTime commented again
func ServerTime(ws *websocket.Conn) {
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
