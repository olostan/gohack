package gohack

import (
	"appengine"
	"appengine/channel"
	"appengine/user"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/api/user", userHandler)
	http.HandleFunc("/api/join", joinHandler)
	http.HandleFunc("/api/message", messageHandler)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	fmt.Fprintf(w, "Hello, %v!", u)
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	tok, err := channel.Create(c, "general")
	if err != nil {
		http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
		c.Errorf("channel.Create: %v", err)
		return
	}
	fmt.Fprintf(w, tok)
}
func readBody(r *http.Request) (string, error) {
	result, err := ioutil.ReadAll(r.Body)
	return string(result), err
}
func messageHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	message, err := readBody(r)
	message = html.EscapeString(message)
	if err != nil {
		http.Error(w, "Couldn't read message", http.StatusInternalServerError)
		c.Errorf("Reading body: %v", err)
		return
	}
	err = channel.SendJSON(c, "general", Message{u.String(), message, time.Now()})
	if err != nil {
		c.Errorf("sending Game: %v", err)
	}
	fmt.Fprintf(w, "Ok")
}

type Message struct {
	Sender string
	Text   string
	When   time.Time
}
