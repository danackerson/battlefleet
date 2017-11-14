package routes

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/danackerson/battlefleet/app"
	"github.com/danackerson/battlefleet/app/structures"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

// CallbackHandler handles login callback from Auth0
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	domain := app.AuthZeroData.Auth0Domain
	conf := &oauth2.Config{
		ClientID:     app.AuthZeroData.Auth0ClientID,
		ClientSecret: app.AuthZeroData.Auth0ClientSecret,
		RedirectURL:  app.AuthZeroData.Auth0CallbackURLString,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "Auth0 token: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "Auth0 client: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err = json.Unmarshal(raw, &profile); err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	session, err := app.SessionStore.Get(r, app.SessionCookieKey)
	if err != nil {
		origError := err.Error()
		if strings.Contains(err.Error(), "no such file or directory") {
			// probably have an old session cookie so recreate
			session, err = app.SessionStore.New(r, app.SessionCookieKey)
			if err != nil {
				if !strings.Contains(err.Error(), "no such file or directory") {
					t, _ := template.New("errorPage").Parse(errorPage)
					t.Execute(w, "recreateSession: "+err.Error()+" (after '"+origError+"')")
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}
			}
		} else {
			t, _ := template.New("errorPage").Parse(errorPage)
			t.Execute(w, "getSession: "+err.Error())
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
	}

	var accountFound *structures.Account
	mongoSession := app.DB.Copy()
	defer mongoSession.Close()
	c := mongoSession.DB("fleetbattle").C("sessions")
	err = c.Find(bson.M{"auth0profile.sub": profile["sub"]}).One(&accountFound)
	if err != nil && err.Error() == "not found" {
		log.Printf("didn't find '%s' in mongoDB\n", profile["name"])
		accountFound = structures.NewAccount(profile["nickname"].(string))
	} else if err != nil && err.Error() != "not found" {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "MongoDB: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	var accountPlaying *structures.Account
	if session.Values[app.AccountKey] != nil {
		accountPlaying = session.Values[app.AccountKey].(*structures.Account)
		if accountPlaying.Commander != app.DefaultCmdrName {
			accountFound.Commander = accountPlaying.Commander
		}
		accountFound.Games = append(accountFound.Games, accountPlaying.Games...)
	}

	accountFound.Auth0Token = token
	accountFound.Auth0Profile = profile
	accountFound.LastLogin = time.Now()

	session.Values[app.AccountKey] = accountFound
	err = session.Save(r, w)
	if err != nil {
		t, _ := template.New("errorPage").Parse(errorPage)
		t.Execute(w, "saveSession: "+err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/account/", http.StatusSeeOther)
}
