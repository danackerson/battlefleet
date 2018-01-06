package routes

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

	go retrieveGame(ws, w, r)
}

func retrieveGame(ws *websocket.Conn, w http.ResponseWriter, r *http.Request) {
	defer ws.Close()

	for {
		// it's helpful to negotiate client side CLOSE requests
		// msg types: https://github.com/gorilla/websocket/blob/master/conn.go#L61
		messageType, msg, err := ws.ReadMessage()
		log.Printf("%d %s (ERR: %v)", messageType, msg, err)

		if messageType != -1 && messageType != 8 {
			session, sessionErr := app.SessionStore.Get(r, app.SessionCookieKey)
			if sessionErr != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("ERR retrieving session:"+sessionErr.Error()))
				return
			}

			if messageType == websocket.PongMessage {
				// while simple it ain't Production ready
				// TODO: https://www.jonathan-petitcolas.com/2015/01/27/playing-with-websockets-in-go.html

				time.Sleep(1 * time.Second)
				ws.WriteMessage(websocket.PingMessage, []byte("PING"))
				continue // skip and block on ws.ReadMessage() above
			}

			account := getAccount(r, session)
			if account != nil {
				game := account.GetGame()

				cmdTermPos := strings.Index(string(msg), ":")

				// switch on battlefleet MSG type (OPEN:, ACK:, MOV: ...)
				switch cmd := string(msg[0:cmdTermPos]); cmd {
				case "MOV": // player has moved a ship - update it
					move := string(msg[cmdTermPos+1:])
					// e36b6892-f076-5f40-aa0a-ba850854d8f3(4,-1)>>(4,-5)
					// this regex is from ws.js => `wsApp.$data.ws.send("MOV:"...`
					re := regexp.MustCompile(`(?P<ShipID>[a-z0-9\-]{36})\((?P<OrigQ>[0-9\-]+)\,(?P<OrigR>[0-9\-]+)\)>>\((?P<DestQ>[0-9\-]+)\,(?P<DestR>[0-9\-]+)\)`)
					match := re.FindStringSubmatch(move)

					// extract ShipID
					shipID := match[1]
					// extract from & to points
					origQ, _ := strconv.ParseFloat(match[2], 16)
					origR, _ := strconv.ParseFloat(match[3], 16)
					destQ, _ := strconv.ParseFloat(match[4], 16)
					destR, _ := strconv.ParseFloat(match[5], 16)
					origin := structures.Point{X: origQ, Y: origR}
					destination := structures.Point{X: destQ, Y: destR}

					// find Ship object in Game and update accordingly
					for _, ship := range game.Ships {
						if ship.ID == shipID && origin == ship.Position {
							//log.Printf("found %v - moving to %v", ship, destination)
							ship.Position = destination
							game.LastTurn = time.Now()
							if e := session.Save(r, w); e != nil {
								t, _ := template.New("errorPage").Parse(errorPage)
								t.Execute(w, "saveSession: "+e.Error())
								http.Redirect(w, r, "/", http.StatusInternalServerError)
								return
							}
							break
						}
					}
				case "OPEN": // initial handshake so send all game info
					/*game.Ships[0].Position = structures.Point{X: 3, Y: -4} // just a test ;)
					game.Ships[0].Class = "luxury"                         // just a test ;)
					game.Ships[0].Type = "fighter"                         // just a test ;)*/
					game.Map = nil // save ~400KB of useless crap on the client side PER WS request
					ws.WriteJSON(game)
				default:
					ws.WriteMessage(websocket.PingMessage, []byte("PING"))
					//ws.WriteMessage(websocket.TextMessage, []byte("Unknown cmd type: "+cmd))
				}
			} else {
				ws.WriteMessage(websocket.TextMessage, []byte("PING"))
			}
		} else {
			log.Printf("Client hung up...good-bye!\n")
			return
		}
	}
}
