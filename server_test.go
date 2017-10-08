package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var host = "http://localhost" + httpPort
var router = mux.NewRouter()
var expiration = time.Now().Add(30 * 24 * time.Hour)

func init() {
	parseEnvVariables()
	setUpMuxHandlers(router)
}

// lot of moving parts in here - no need to rewrite 25 times for testing
// be careful of sessions and cookies
func prepareServeHTTP(requestType string, requestURL string, sessionCookie string, session *sessions.Session) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(requestType, requestURL, nil)
	res := httptest.NewRecorder()

	if sessionCookie != "" {
		req.AddCookie(&http.Cookie{Name: sessionCookieKey, Value: sessionCookie, Expires: expiration})
	}

	if session == nil {
		session, _ = store.Get(req, sessionCookieKey)
		session.Values[gameUUIDKey] = newGameUUID
	} else {
		log.Printf("gameUUID in sessionPrep: %s", session.Values[gameUUIDKey])
	}
	if e := session.Save(req, res); e != nil {
		panic(e) // for now
	}

	return res, req
}

// TestHome1Game1Home2
// 1. visits homepage with no existing session
// 2. creates and visits new game page
// 3. revisits homepage with existing session
// 4. create a new game with existing session
// 5. goto new game page and verify *new* game :)
func TestHome1Game1Home2(t *testing.T) {
	t.Parallel()

	sessionCookie := ""

	// FIRST homepage visit - no session
	res, req := prepareServeHTTP("GET", host+"/", sessionCookie, nil)
	router.ServeHTTP(res, req)
	session, sessErr := store.Get(req, sessionCookieKey)
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
	res, req = prepareServeHTTP("GET", host+exp1, sessionCookie, session)
	router.ServeHTTP(res, req)
	session, _ = store.Get(req, sessionCookieKey)
	newGameURL := res.Header().Get("Location")
	gameUUID := session.Values[gameUUIDKey].(string)
	if res.Code != 301 || !strings.Contains(newGameURL, gameUUID) {
		t.Fatalf("Expecting redirect to /games/%s", gameUUID)
	}

	// THIRD follow redirect to The Game page ;)
	res, req = prepareServeHTTP("GET", host+newGameURL, sessionCookie, session)
	router.ServeHTTP(res, req)
	act = res.Body.String()
	if res.Code != 200 || !strings.Contains(act, "Engage!") {
		log.Printf("%s\n", act)
		t.Fatalf("Expected working game page but didn't get it (%d)", res.Code)
	}

	// FOURTH revisit Home and verify the "Rejoin Fleet!" @ gameUUID URL
	res, req = prepareServeHTTP("GET", host+"/", sessionCookie, session)
	router.ServeHTTP(res, req)
	session, sessErr = store.Get(req, sessionCookieKey)
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
	res, req = prepareServeHTTP("GET", host+"/games/"+newGameUUID, sessionCookie, session)
	router.ServeHTTP(res, req)
	session, _ = store.Get(req, sessionCookieKey)
	newGameURL = res.Header().Get("Location")
	if res.Code != 301 && strings.Contains(newGameURL, gameUUID) {
		t.Fatalf("Expecting redirect to %s", newGameURL)
	} else {
		log.Printf("redirect to %s", newGameURL)
	}

	// SIXTH visit new game page and verify new gameUUID
	res, req = prepareServeHTTP("GET", host+newGameURL, sessionCookie, session)
	router.ServeHTTP(res, req)
	session, _ = store.Get(req, sessionCookieKey)
	gameUUIDNew := session.Values[gameUUIDKey].(string)
	if gameUUIDNew == gameUUID {
		t.Fatalf("Expected new gameUUID (%s)", gameUUIDNew)
	}
	act = res.Body.String()
	if res.Code != 200 || !strings.Contains(act, "Engage!") {
		log.Printf("%s\n", act)
		t.Fatalf("Expected working game page but didn't get it (%d)", res.Code)
	}
}

func TestWebSocketConnect(t *testing.T) {
	// TODO
	t.Parallel()
}

func TestVersion(t *testing.T) {
	t.Parallel()

	res, req := prepareServeHTTP("POST", "http://localhost"+httpPort+"/version", "", nil)
	router.ServeHTTP(res, req)

	// Verify we get a TEST build version id back
	exp := "{\"build\":\"TEST\",\"version\":\"https://circleci.com/gh/danackerson/battlefleet/TEST\"}"
	act := res.Body.String()
	if !strings.Contains(act, exp) {
		t.Fatalf("Expected %s got %s", exp, act)
	}
}
