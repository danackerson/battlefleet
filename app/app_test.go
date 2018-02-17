package app_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/routes"
	"github.com/danackerson/battlefleet/app/structures"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var router *mux.Router

func init() {
	isMainExec := false
	app.Init(isMainExec)
	os.Setenv("TEMPLATE_DIR", "templates") // override for testing
	router = routes.SetUpMuxHandlers(isMainExec)
}

type testRequestContext struct {
	requestType   string
	requestURL    string
	sessionCookie string
	headers       map[string]string
	session       *sessions.Session
	formVariables *strings.Reader
}

func prepareServeHTTP(context *testRequestContext) (*httptest.ResponseRecorder, *http.Request) {
	host := "http://localhost" + app.HTTPPort
	req := httptest.NewRequest(context.requestType, host+context.requestURL, nil)

	if context.formVariables != nil {
		req = httptest.NewRequest(context.requestType, host+context.requestURL, context.formVariables)
	}
	res := httptest.NewRecorder()

	if context.requestType == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if context.sessionCookie != "" {
		req.AddCookie(&http.Cookie{
			Name:    app.SessionCookieKey,
			Value:   context.sessionCookie,
			Expires: time.Now().Add(30 * 24 * time.Hour),
		})
	}
	if len(context.headers) > 0 {
		for key, value := range context.headers {
			req.Header.Add(key, value)
		}
	}

	if context.session == nil {
		context.session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
		context.session.Values[app.GameUUIDKey] = structures.NewGameUUID
	}

	if e := context.session.Save(req, res); e != nil {
		panic(e) // for now
	}

	return res, req
}

// TestNewAccountNewGameDeleteGameDeleteAccount
func TestNewAccountNewGameDeleteGameDeleteAccount(t *testing.T) {
	t.Parallel()

	sessionCookie, session := testFirstHomePageVisit(t)
	req, resp := testNewGame(t, sessionCookie, session)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	account := session.Values[app.AccountKey].(*structures.Account)
	newGameURL := resp.Header().Get("Location")
	session = testVisitGamePage(t, newGameURL, sessionCookie, session)

	// FOURTH: visit the account page
	context := &testRequestContext{
		requestType: "GET", requestURL: "/account/",
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req)

	gameUUID := session.Values[app.GameUUIDKey].(string)
	act := res.Body.String()
	exp1 := "confirmGameDeletion(&#34;" + gameUUID + "&#34;);"
	exp2 := "Welcome, " + account.Commander
	if !strings.Contains(act, exp1) && !strings.Contains(act, exp2) {
		t.Fatalf("Expected %s\ngot %s", exp1, act)
	}

	// FIFTH: delete the game
	context = &testRequestContext{
		requestType: "GET", requestURL: "/games/" + gameUUID + "?action=delete",
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	account = session.Values[app.AccountKey].(*structures.Account)
	if len(account.Games) != 0 &&
		account.CurrentGameID != structures.NewGameUUID &&
		session.Values[app.GameUUIDKey] != "" {
		t.Fatalf("Expected 0 Games for current account")
	}

	// SIXTH: delete the account
	context = &testRequestContext{
		requestType: "GET", requestURL: "/account/?action=delete",
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	newAccount := session.Values[app.AccountKey]
	cookieString := res.Header().Get("Set-Cookie")
	cookieString = strings.TrimLeft(cookieString, app.SessionCookieKey+"=")
	newSessionCookie := strings.Split(cookieString, ";")[0]

	if newAccount != nil && sessionCookie != newSessionCookie {
		t.Fatalf("Expected a brand new session with No Account!")
	}
}

func testFirstHomePageVisit(t *testing.T) (string, *sessions.Session) {
	// FIRST homepage visit - no session
	context := &testRequestContext{
		requestType: "GET", requestURL: "/",
		sessionCookie: "", session: nil, formVariables: nil,
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req) // new session appears right here!
	session, sessErr := app.SessionStore.Get(req, app.SessionCookieKey)
	if sessErr != nil {
		log.Printf("sessErr 1: %v\n", sessErr)
	}
	cookieString := res.Header().Get("Set-Cookie")
	cookieString = strings.TrimLeft(cookieString, app.SessionCookieKey+"=")
	sessionCookie := strings.Split(cookieString, ";")[0]

	// Verify we see a URL with the "__new__" as "New game!"
	// and NO button to rejoin the fleet...
	exp1 := "/games/" + session.Values[app.GameUUIDKey].(string)
	exp2 := "New game!"
	exp3 := "Rejoin fleet"
	act := res.Body.String()
	if !strings.Contains(act, exp1) || !strings.Contains(act, exp2) || strings.Contains(act, exp3) {
		t.Fatalf("Expected %s\ngot %s", exp1, act)
	}

	return sessionCookie, session
}

func testNewGame(t *testing.T, sessionCookie string, session *sessions.Session) (req *http.Request, resp *httptest.ResponseRecorder) {
	context := &testRequestContext{
		requestType:   "POST",
		requestURL:    "/games/" + session.Values[app.GameUUIDKey].(string),
		sessionCookie: sessionCookie, session: session,
		formVariables: strings.NewReader("cmdrName=Shade"),
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	gameUUID := session.Values[app.GameUUIDKey].(string)
	account := session.Values[app.AccountKey].(*structures.Account)
	cmdrNamed := account.Commander
	newGameURL := res.Header().Get("Location")
	if res.Code != 301 || cmdrNamed != "Shade" || !strings.Contains(newGameURL, gameUUID) {
		t.Fatalf("Expecting redirect to /games/%s", gameUUID)
	}

	return req, res
}

func testVisitGamePage(t *testing.T, newGameURL string,
	sessionCookie string, session *sessions.Session) *sessions.Session {
	// THIRD follow redirect to The Game page ;)
	context := &testRequestContext{
		requestType:   "GET",
		requestURL:    newGameURL,
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	account := session.Values[app.AccountKey].(*structures.Account)
	if len(account.Games) != 1 {
		t.Fatalf("Expected only one game!")
	}
	act := res.Body.String()
	if res.Code != 200 || !strings.Contains(act, "Engage!") {
		log.Printf("%s\n", act)
		t.Fatalf("Expected working game page but didn't get it (%d)", res.Code)
	}

	return session
}

// TestHome1NewGame1Home2NewGame2
// 1. visits homepage with no existing session
// 2. creates and visits new game page
// 3. revisits homepage with existing session
// 4. create a new game with existing session
// 5. goto new game page and verify *new* game :)
func TestHome1Game1Home2(t *testing.T) {
	t.Parallel()

	sessionCookie, session := testFirstHomePageVisit(t)

	req, resp := testNewGame(t, sessionCookie, session)
	newGameURL := resp.Header().Get("Location")
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	gameUUID := session.Values[app.GameUUIDKey].(string)
	session = testVisitGamePage(t, newGameURL, sessionCookie, session)

	// FOURTH revisit Home and verify the "Rejoin Fleet!" @ gameUUID URL
	context := &testRequestContext{
		requestType: "GET", requestURL: "/",
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, sessErr := app.SessionStore.Get(req, app.SessionCookieKey)
	account := session.Values[app.AccountKey].(*structures.Account)
	if len(account.Games) != 1 {
		t.Fatalf("Expected no new games to be created!")
	}
	exp1 := ""
	exp3 := "Rejoin fleet"
	if sessErr != nil {
		log.Printf("sessErr 3: %v\n", sessErr)
	} else {
		exp1 = "/games/" + session.Values[app.GameUUIDKey].(string)
	}
	act := res.Body.String()
	if !strings.Contains(act, exp1) && !strings.Contains(act, exp3) {
		t.Fatalf("Expected %s\ngot %s", exp1, act)
	}

	// FIFTH visit "/games/__new__" and verify a new gameUUIDKey
	context = &testRequestContext{
		requestType:   "GET",
		requestURL:    "/games/" + structures.NewGameUUID,
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	newGameURL = res.Header().Get("Location")
	if res.Code != 301 && strings.Contains(newGameURL, gameUUID) {
		t.Fatalf("Expecting redirect to %s", newGameURL)
	}

	// SIXTH visit new game page and verify new gameUUID
	context = &testRequestContext{
		requestType:   "POST",
		requestURL:    newGameURL,
		sessionCookie: sessionCookie, session: session,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	gameUUIDNew := session.Values[app.GameUUIDKey].(string)
	account = session.Values[app.AccountKey].(*structures.Account)
	if gameUUIDNew == gameUUID || len(account.Games) != 2 {
		t.Fatalf("Expected new gameUUID (%s) but got %s (and %d games?)", gameUUIDNew, gameUUID, len(account.Games))
	}
	act = res.Body.String()
	if res.Code != 200 || !strings.Contains(act, "Engage!") {
		log.Printf("%s\n", act)
		t.Fatalf("Expected working game page but didn't get it (%d)", res.Code)
	}

	// delete session as teardown
	session.Options.MaxAge = -1
	if e := session.Save(req, res); e != nil {
		panic(e) // for now
	}
}

// TestNonLoggedInWebSocketConnect currently just verifies working WebSocket handler
func TestNonLoggedInWebSocketConnect(t *testing.T) {
	t.Parallel()

	// 1. open websocket and get response for user without an Account
	srv := httptest.NewServer(http.HandlerFunc(routes.ServeWebSocket))
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	//log.Printf("Testing server @ %+v\n", u)

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("cannot make websocket connection: %v => %v", err, resp)
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte("OPEN"))
	if err != nil {
		t.Fatalf("cannot write message: %v", err)
	}
	msgType, msgText, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	} else if msgType != 1 {
		t.Fatalf("unexpected msgType (%d) => %s", msgType, msgText)
	} else if string(msgText) != "PING" {
		t.Fatalf("unexpected msgText %s (expected 'PING')", msgText)
	}

	// 2. close websocket and verify graceful shutdown
	err = conn.Close()
	if err != nil {
		t.Fatalf("cannot close WebSocket: %v", err)
	}
}

// TestLoggedInWebSocketConnect currently just verifies working WebSocket handler
func TestLoggedInWebSocketConnect(t *testing.T) {
	t.Parallel()

	// 1. createAccount/Game via test methods above & get real values here
	sessionCookie, session := testFirstHomePageVisit(t)
	req, resp := testNewGame(t, sessionCookie, session)
	session, _ = app.SessionStore.Get(req, app.SessionCookieKey)
	account := session.Values[app.AccountKey].(*structures.Account)
	game := account.GetGame()
	shipID := game.Ships[0].ID

	session.Save(req, resp) // persist to disk so WebSocket server can reuse

	// 2. open websocket and get server time
	//var header map[string][]string
	srv := httptest.NewServer(http.HandlerFunc(routes.ServeWebSocket))
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	//log.Printf("Testing server @ %+v\n", u)

	requestHeader := &http.Header{}
	requestHeader.Add("Origin", "http://"+u.Host)
	requestHeader.Add("Cookie", resp.Header().Get("Set-Cookie"))

	conn, wsResp, err := websocket.DefaultDialer.Dial(u.String(), *requestHeader)
	if err != nil {
		t.Fatalf("cannot make websocket connection: %v => %v", err, wsResp)
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"OPEN"}`))
	if err != nil {
		t.Fatalf("cannot write message: %v", err)
	}
	msgType, msgText, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	} else if msgType != 1 {
		t.Fatalf("unexpected msgType (%d) => %s", msgType, msgText)
	} else if !strings.Contains(string(msgText), shipID) {
		t.Fatalf("unexpected msgText %s (expected 'PONG')", msgText)
	}

	// perform a MOV cmd and verify ship has been moved on server side
	destQ := -1
	destR := -2
	movCmd := []byte(`{"cmd":"MOV","payload":{"ID":"` + shipID +
		`","origin":{"Q":0,"R":0},"destination":{"Q":` +
		strconv.Itoa(destQ) + `,"R":` + strconv.Itoa(destR) + `}}}`)
	err = conn.WriteMessage(websocket.TextMessage, movCmd)
	if err != nil {
		t.Fatalf("cannot write message: %v", err)
	}
	msgType, msgText, err = conn.ReadMessage()

	// parse game JSON response
	var serverResponse map[string]interface{}
	jsonError := json.Unmarshal(msgText, &serverResponse)
	if jsonError != nil {
		log.Printf("server disconnect -> ERR unmarshalling JSON response: %s\n", jsonError.Error())
		return
	}

	saveTime := strings.Split(serverResponse["LastTurn"].(string), "+")[0]
	// datetime is saved with 5 OR 6 digits of millisecond precision :(
	saveTimeNoMillis := strings.Split(saveTime, ".")[0]
	saveTimeMillis := strings.Split(saveTime, ".")[1]
	saveTime = saveTimeNoMillis + "." + saveTimeMillis[0:5]

	layout := "2006-01-02T15:04:05.00000"
	modifyGameDate, terr := time.Parse(layout, saveTime)
	if terr != nil {
		log.Println(terr.Error())
	}

	// verify new LastTurn value is *after* original `game.LastTurn`
	if !modifyGameDate.After(game.LastTurn) {
		t.Fatalf("modify game date (%s) NOT after create date (%s)", modifyGameDate, game.LastTurn)
	}

	modShipsArr := serverResponse["Ships"].([]interface{})
	modShipObj := modShipsArr[0].(map[string]interface{})
	if modShipObj["ID"] != shipID {
		t.Fatalf("modified ship (%s) is not same as original ship (%s)", modShipObj["ID"], shipID)
	}
	modShipPos := modShipObj["Position"].(map[string]interface{})
	var modShipPosX = modShipPos["X"].(float64)
	var modShipPosY = modShipPos["Y"].(float64)
	if int(modShipPosX) != destQ && int(modShipPosY) != destR {
		t.Fatalf("new ship loc (%d, %d) is same as init ship loc (0, 0)", int(modShipPosX), int(modShipPosY))
	}

	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	} else if msgType != 1 {
		t.Fatalf("unexpected msgType (%d) => %s", msgType, msgText)
	} else if !strings.Contains(string(msgText), shipID) {
		t.Fatalf("unexpected msgText %s (expected 'PONG')", msgText)
	}

	// 3. close websocket and verify graceful shutdown
	err = conn.Close()
	if err != nil {
		t.Fatalf("cannot close WebSocket: %v", err)
	}

	// delete session as teardown
	session.Options.MaxAge = -1
	if e := session.Save(req, resp); e != nil {
		panic(e) // for now
	}
}

// TODO: add this to the footer of all pages!
func TestVersion(t *testing.T) {
	t.Parallel()

	context := &testRequestContext{
		requestType:   "POST",
		requestURL:    "/version",
		sessionCookie: "cook", session: nil, formVariables: nil,
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req)

	// Verify we get the build version id back
	exp := "{\"build\":\"" + app.Version + "\",\"version\":\"https://circleci.com/gh/danackerson/battlefleet/" + app.Version + "\"}"
	act := res.Body.String()
	if !strings.Contains(act, exp) {
		t.Fatalf("Expected %s got %s", exp, act)
	}

	context.session.Options.MaxAge = -1
	if e := context.session.Save(req, res); e != nil {
		panic(e) // for now
	}
}
