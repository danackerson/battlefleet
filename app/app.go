package app

import (
	"encoding/gob"
	"errors"
	"html/template"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/danackerson/battlefleet/app/structures"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	mgo "gopkg.in/mgo.v2"
)

var (
	ProdSession      = false
	HTTPPort         = ":8083"
	SessionCookieKey = "battlefleetID"
	AccountIDKey     = "ownerAccountID"
	AccountKey       = "ownerAccount"
	CmdrNameKey      = "cmdrName"
	DefaultCmdrName  = "stranger!"
	GameUUIDKey      = "gameUUID"
	URIScheme        string
	Version          string
	DB               *mgo.Session
	SessionStore     *sessions.FilesystemStore
	SessionMaxAge    int
	AuthZeroData     Auth0Data
)

var mongoDBUser string
var mongoDBPass string
var mongoDBHost string
var mongoDBName = "fleetbattle"
var mongoDBCollection = "sessions"

// Auth0Data for storing Auth0 variables
type Auth0Data struct {
	Auth0ClientID          string
	Auth0ClientSecret      string
	Auth0Domain            string
	Auth0CallbackURLString string
	Auth0CallbackURL       template.URL
}

// Init all the state and session information for the application
func Init(isMainExec bool) {
	prepareSessionEnvironment(isMainExec)
	setupMongoDBSession()
}

// GetOutboundIP or preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func prepareSessionEnvironment(isMainExec bool) {
	SessionStore = sessions.NewFilesystemStore("/tmp", []byte(os.Getenv("bfSecret")))
	SessionStore.MaxLength(32 * 1024) // else securecookie: value too long error
	gob.Register(&structures.Account{})
	gob.Register(&structures.Game{})
	gob.Register(&structures.Ship{})

	ProdSession, _ = strconv.ParseBool(os.Getenv("prodSession"))
	if !ProdSession {
		URIScheme = "http"
		SessionMaxAge := 7 * 24 * 3600 // 1 week

		SessionStore.Options = &sessions.Options{
			Path:   "/",
			Domain: "localhost",   // NOTE: use `GetOutboundIP().String()` for local, Mobile device testing
			MaxAge: SessionMaxAge, // one week
		}

		// load test vars from .env
		envDir := "../.env"
		if isMainExec {
			envDir = ".env"
		}
		err := godotenv.Load(envDir)
		if err != nil {
			log.Fatalf("env var load failure: %s", err.Error())
		}
	} else {
		URIScheme = "https"
		SessionMaxAge := 3600 * 24 * 365 // 1 year expiration
		SessionStore.Options = &sessions.Options{
			Path:     "/",
			Domain:   "battlefleet.eu",
			MaxAge:   SessionMaxAge,
			Secure:   true,
			HttpOnly: true,
		}
	}

	AuthZeroData = Auth0Data{
		os.Getenv("AUTH0_CLIENT_ID"),
		os.Getenv("AUTH0_CLIENT_SECRET"),
		os.Getenv("AUTH0_DOMAIN"),
		os.Getenv("AUTH0_CALLBACK_URL"),
		template.URL(os.Getenv("AUTH0_CALLBACK_URL")),
	}

	Version = os.Getenv("CIRCLE_BUILD_NUM")
	if Version == "" {
		Version = "TEST"
	}
}

func setupMongoDBSession() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("mongoDBHost")},
		Timeout:  10 * time.Second,
		Database: os.Getenv("mongoDBName"),
		Username: os.Getenv("mongoDBUser"),
		Password: os.Getenv("mongoDBPass"),
	}

	err := errors.New("")
	DB, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("No connection to MongoDB: %s", err.Error())
	}
	DB.SetMode(mgo.Monotonic, true)
}
