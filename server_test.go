package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/danackerson/battlefleet/structures"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var host = "http://localhost" + httpPort
var wsHost = "ws://localhost" + httpPort
var router = mux.NewRouter()

func init() {
	prepareSessionEnvironment()
	setUpMuxHandlers(router)
	setUpFuncMaps()
}

type testRequestContext struct {
	requestType   string
	requestURL    string
	sessionCookie string
	session       *sessions.Session
	formVariables *strings.Reader
}

func prepareServeHTTP(context *testRequestContext) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(context.requestType, context.requestURL, nil)
	if context.formVariables != nil {
		req = httptest.NewRequest(context.requestType, context.requestURL, context.formVariables)
	}
	res := httptest.NewRecorder()

	if context.requestType == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if context.sessionCookie != "" {
		req.AddCookie(&http.Cookie{Name: sessionCookieKey, Value: context.sessionCookie, Expires: time.Now().Add(30 * 24 * time.Hour)})
	}

	if context.session == nil {
		context.session, _ = sessionStore.Get(req, sessionCookieKey)
		context.session.Values[gameUUIDKey] = structures.NewGameUUID
	}

	if e := context.session.Save(req, res); e != nil {
		panic(e) // for now
	}

	return res, req
}

// TODO - remember to cleanup session afterwards!
// TestNewAccountNewGameDeleteGameDeleteAccount

// TODO - remember to cleanup session afterwards!
// TestNewAccountNewGame1NewGame2NewGame3NoGame4DeleteGame3HomePageNewGame

// TestHome1NewGame1Home2NewGame2
// 1. visits homepage with no existing session
// 2. creates and visits new game page
// 3. revisits homepage with existing session
// 4. create a new game with existing session
// 5. goto new game page and verify *new* game :)
func TestHome1Game1Home2(t *testing.T) {
	t.Parallel()

	sessionCookie := "" // comment

	// FIRST homepage visit - no session
	context := &testRequestContext{
		requestType: "GET", requestURL: host + "/",
		sessionCookie: sessionCookie, session: nil, formVariables: nil,
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req) // new session appears right here!
	session, sessErr := sessionStore.Get(req, sessionCookieKey)
	if sessErr != nil {
		log.Printf("sessErr 1: %v\n", sessErr)
	}
	cookieString := res.Header().Get("Set-Cookie")
	cookieString = strings.TrimLeft(cookieString, sessionCookieKey+"=")
	sessionCookie = strings.Split(cookieString, ";")[0]

	// Verify we see a URL with the "__new__" as "New game!"
	// and NO button to rejoin the fleet...
	exp1 := "/games/" + session.Values[gameUUIDKey].(string)
	exp2 := "New game!"
	exp3 := "Rejoin fleet"
	act := res.Body.String()
	if !strings.Contains(act, exp1) || !strings.Contains(act, exp2) || strings.Contains(act, exp3) {
		t.Fatalf("Expected %s\ngot %s", exp1, act)
	}

	// SECOND visit New Game page
	context = &testRequestContext{
		requestType:   "POST",
		requestURL:    host + exp1,
		sessionCookie: sessionCookie, session: session,
		formVariables: strings.NewReader("cmdrName=Shade"),
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = sessionStore.Get(req, sessionCookieKey)
	gameUUID := session.Values[gameUUIDKey].(string)
	account := session.Values[accountKey].(*structures.Account)
	cmdrNamed := account.Commander
	newGameURL := res.Header().Get("Location")
	if res.Code != 301 || cmdrNamed != "Shade" || !strings.Contains(newGameURL, gameUUID) {
		t.Fatalf("Expecting redirect to /games/%s", gameUUID)
	}

	// THIRD follow redirect to The Game page ;)
	context = &testRequestContext{
		requestType:   "GET",
		requestURL:    host + newGameURL,
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = sessionStore.Get(req, sessionCookieKey)
	account = session.Values[accountKey].(*structures.Account)
	if len(account.Games) != 1 {
		t.Fatalf("Expected only one game!")
	}
	act = res.Body.String()
	if res.Code != 200 || !strings.Contains(act, "Engage!") {
		log.Printf("%s\n", act)
		t.Fatalf("Expected working game page but didn't get it (%d)", res.Code)
	}

	// FOURTH revisit Home and verify the "Rejoin Fleet!" @ gameUUID URL
	context = &testRequestContext{
		requestType: "GET", requestURL: host + "/",
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, sessErr = sessionStore.Get(req, sessionCookieKey)
	account = session.Values[accountKey].(*structures.Account)
	if len(account.Games) != 1 {
		t.Fatalf("Expected no new games to be created!")
	}
	if sessErr != nil {
		log.Printf("sessErr 3: %v\n", sessErr)
	} else {
		exp1 = "/games/" + session.Values[gameUUIDKey].(string)
	}
	act = res.Body.String()
	if !strings.Contains(act, exp1) && !strings.Contains(act, exp3) {
		t.Fatalf("Expected %s\ngot %s", exp1, act)
	}

	// FIFTH visit "/games/__new__" and verify a new gameUUIDKey
	context = &testRequestContext{
		requestType:   "GET",
		requestURL:    host + "/games/" + structures.NewGameUUID,
		sessionCookie: sessionCookie, session: session, formVariables: nil,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = sessionStore.Get(req, sessionCookieKey)
	newGameURL = res.Header().Get("Location")
	if res.Code != 301 && strings.Contains(newGameURL, gameUUID) {
		t.Fatalf("Expecting redirect to %s", newGameURL)
	}

	// SIXTH visit new game page and verify new gameUUID
	context = &testRequestContext{
		requestType:   "POST",
		requestURL:    host + newGameURL,
		sessionCookie: sessionCookie, session: session,
	}
	res, req = prepareServeHTTP(context)
	router.ServeHTTP(res, req)
	session, _ = sessionStore.Get(req, sessionCookieKey)
	gameUUIDNew := session.Values[gameUUIDKey].(string)
	account = session.Values[accountKey].(*structures.Account)
	if gameUUIDNew == gameUUID || len(account.Games) != 2 {
		t.Fatalf("Expected new gameUUID (%s)", gameUUIDNew)
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

// TestWebSocketConnect currently just verifies working WebSocket handler (datetime)
// TODO: verify valid session w/ gameUUID before responding
func TestWebSocketConnect(t *testing.T) {
	t.Parallel()

	/*sessionCookie := ""

	// 0. visit "/games/__new__" and get a new gameUUIDKey
	res, req := prepareServeHTTP("GET", host+"/games/"+newGameUUID, sessionCookie, nil)
	router.ServeHTTP(res, req)
	session, _ := store.Get(req, sessionCookieKey)
	newGameURL := res.Header().Get("Location")

	// 1. visit new game page and verify new gameUUID
	reader := strings.NewReader("cmdrName=WSUser")
	res, req = prepareServeHTTP("POST", host+newGameURL, sessionCookie, session, reader)
	router.ServeHTTP(res, req)
	//session, _ = store.Get(req, sessionCookieKey)
	//gameUUID := session.Values[gameUUIDKey].(string)*/

	// 2. open websocket and get server time
	srv := httptest.NewServer(http.HandlerFunc(serveWebSocket))
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
	}

	// 3. close websocket and verify graceful shutdown
	err = conn.Close()
	if err != nil {
		t.Fatalf("cannot close WebSocket: %v", err)
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()

	context := &testRequestContext{
		requestType:   "POST",
		requestURL:    host + "/version",
		sessionCookie: "cook", session: nil, formVariables: nil,
	}
	res, req := prepareServeHTTP(context)
	router.ServeHTTP(res, req)

	// Verify we get the build version id back
	exp := "{\"build\":\"" + version + "\",\"version\":\"https://circleci.com/gh/danackerson/battlefleet/" + version + "\"}"
	act := res.Body.String()
	if !strings.Contains(act, exp) {
		t.Fatalf("Expected %s got %s", exp, act)
	}

	context.session.Options.MaxAge = -1
	if e := context.session.Save(req, res); e != nil {
		panic(e) // for now
	}
}
