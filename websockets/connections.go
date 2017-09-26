package websockets

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// ServerTime commented again
func ServerTime(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()

	for {
		// even though we don't do much with the Reads (yet)
		// it's helpful to negotiate client side CLOSE requests
		messageType, msg, err := ws.ReadMessage()
		if err == nil {
			//log.Printf("%d: %s", messageType, msg)
			if messageType == 2 {
			} else if len(msg) == 3 {
			}
		} else {
			log.Printf("%v\n", err)
		}
		serverTimeBytes := time.Now().Format(time.UnixDate)
		if err := ws.WriteMessage(websocket.TextMessage, []byte(serverTimeBytes)); err != nil {
			log.Printf("txtMsg: " + err.Error())
			return
		}
		time.Sleep(1 * time.Second)
	}

}
