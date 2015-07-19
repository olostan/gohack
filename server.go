package gohack

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
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

type InitialReply struct {
	User     string
	Messages []Message
}

func reverseMessages(a []Message) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	u := user.Current(c)
	initial := InitialReply{u.String(), make([]Message, 0, 10)}

	q := datastore.NewQuery("Message").Ancestor(channelKey(c)).Order("-When").Limit(10)
	if _, err := q.GetAll(c, &initial.Messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reverseMessages(initial.Messages)
	res, err := json.Marshal(initial)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(res))
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
	messageText, err := readBody(r)
	messageText = html.EscapeString(messageText)
	if err != nil {
		http.Error(w, "Couldn't read message", http.StatusInternalServerError)
		c.Errorf("Reading body: %v", err)
		return
	}
	message := Message{u.String(), messageText, time.Now()}
	err = channel.SendJSON(c, "general", message)
	if err != nil {
		c.Errorf("sending Game: %v", err)
	}
	key := datastore.NewIncompleteKey(c, "Message", channelKey(c))
	_, err = datastore.Put(c, key, &message)
	if err != nil {
		c.Errorf("Storage message: %v", err)
	}
	fmt.Fprintf(w, "Ok")
}

func channelKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "GoHack", "general", 0, nil)
}

type Message struct {
	Sender string
	Text   string
	When   time.Time
}
