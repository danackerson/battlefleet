package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/structures"
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

	go serverTime(ws)
	//go retrieveGame(ws, w, r)
}

type serverTimeStr struct {
	Action    string `json:"action"`
	Time      int64  `json:"time,omitempty"`
	Mutation  string `json:"mutation,omitempty"`
	SessionID string `json:"sessionID,omitempty"`
}

func serverTime(ws *websocket.Conn) {
	defer ws.Close()

	for {
		// msg types: https://github.com/gorilla/websocket/blob/master/conn.go#L61

		// Realize if we block here with a ReadJSON method, the server won't forever WriteJSON ;)
		var clientMsg interface{}
		err := ws.ReadJSON(&clientMsg)
		if err != nil {
			log.Println(err.Error())
		}
		log.Printf("RCVD %v", clientMsg)
		// TODO: read clientMsg and verify Session

		srvTime := serverTimeStr{}
		srvTime.Action = "currentServerTime"
		srvTime.Time = time.Now().Unix()
		//srvTime.Mutation = "SOCKET_ONMESSAGE"
		if err := ws.WriteJSON(srvTime); err != nil {
			log.Printf("SENT %v (ERR: %s)", srvTime, err.Error())
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func startPing(ws *websocket.Conn) *time.Ticker {
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		select {
		case <-ticker.C:
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("unable to PING: %s\n", err.Error())
				return
			}
		}
	}()

	return ticker
}

func retrieveGame(ws *websocket.Conn, w http.ResponseWriter, r *http.Request) {
	ticker := startPing(ws)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()

	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil && messageType != -1 {
			log.Println("server disconnect -> ERR reading WS msg: " + err.Error())
			return
		}
		log.Printf("%d %s (ERR: %v)", messageType, msg, err)

		if messageType != -1 && messageType != 8 {
			session, sessionErr := app.SessionStore.Get(r, app.SessionCookieKey)
			if sessionErr != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("ERR retrieving session: "+sessionErr.Error()))
				log.Println("server disconnect -> ERR reading WS session: " + sessionErr.Error())
				return
			}

			if messageType == websocket.PongMessage {
				log.Printf("PONG\n")
				continue // skip rest and block on ws.ReadMessage() above
			}

			account := getAccount(session, "UNKNOWN")
			if account != nil {
				game := account.GetGame()

				var clientRequest map[string]interface{}
				jsonError := json.Unmarshal(msg, &clientRequest)
				if jsonError != nil {
					log.Printf("server disconnect -> ERR unmarshalling JSON request: %s\n", jsonError.Error())
					return
				}

				// switch on battlefleet MSG type (OPEN:, ACK:, MOV: ...)
				cmd := clientRequest["cmd"].(string)
				switch cmd {
				case "MOV": // player has moved a ship - update it
					move := clientRequest["payload"].(map[string]interface{})
					shipID := move["ID"]
					originStr := move["origin"].(map[string]interface{})
					destinationStr := move["destination"].(map[string]interface{})
					// extract from & to points
					origQ, _ := originStr["Q"]
					origR, _ := originStr["R"]
					destQ, _ := destinationStr["Q"]
					destR, _ := destinationStr["R"]
					origin := structures.Point{X: origQ.(float64), Y: origR.(float64)}
					destination := structures.Point{X: destQ.(float64), Y: destR.(float64)}

					// find Ship object in Game and update accordingly
					for _, ship := range game.Ships {
						if ship.ID == shipID && origin == ship.Position {
							//log.Printf("found %v - moving to %v", ship, destination)
							ship.Position = destination
							game.LastTurn = time.Now()
							if e := session.Save(r, w); e != nil {
								log.Println("server disconnecting -> MOV ship: sessionSave failed: " + e.Error())
								return
							}
							break
						}
					}

					writeErr := ws.WriteJSON(game)
					if writeErr != nil {
						log.Println("server disconnecting -> MOV writeJSON response failed: " + writeErr.Error())
						return
					}
				case "OPEN": // initial handshake so send all game info
					/*game.Ships[0].Position = structures.Point{X: 3, Y: -4} // just a test ;)
					game.Ships[0].Class = "luxury"                         // just a test ;)
					game.Ships[0].Type = "fighter"                         // just a test ;)*/
					game.Map = nil // save ~400KB of useless crap on the client side PER WS request
					writeErr := ws.WriteJSON(game)
					if writeErr != nil {
						log.Println("server disconnecting -> OPEN writeJSON response failed: " + writeErr.Error())
						return
					}
				default:
					writeErr := ws.WriteMessage(websocket.TextMessage, []byte("Unknown cmd type: "+cmd))
					if writeErr != nil {
						log.Println("server disconnecting -> unknown cmd param passed to websocket.ReadMessage(): " + writeErr.Error())
						return
					}
				}
			} else {
				writeErr := ws.WriteMessage(websocket.TextMessage, []byte("PING")) // used for unit test
				if writeErr != nil {
					log.Println("server disconnecting -> OPEN writeJSON response failed: " + writeErr.Error())
					return
				}
			}
		} else {
			log.Printf("Client hung up...good-bye!\n")
			return
		}
	}
}
