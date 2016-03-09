package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	JQuery template.JS
	Css    template.CSS
	Data   string
	Logo   string
	Port   string
}

func NewPage() (*Page, error) {
	p := &Page{
		Port: Port,
		Data: "{}",
	}

	if data := registry.Json(); data != "null" {
		p.Data = data
	}

	l, err := Asset("assets/beepboop_2x.ca5eff.png")
	if err != nil {
		log.Printf("Error %s", err.Error())
	}
	p.Logo = base64.StdEncoding.EncodeToString(l)

	return p, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p, err := NewPage()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, p)
}

var indexTmpl, _ = Asset("assets/index.html")
var t = template.Must(template.New("index.html").Parse(string(indexTmpl)))

func renderTemplate(w http.ResponseWriter, p *Page) {
	err := t.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	css, err := Asset("assets/pure-min.css")
	if err != nil {
		log.Printf("Error %s", err.Error())
	}

	w.Header().Set("Content-Type", "text/css")
	cssStr := string(css)
	fmt.Fprint(w, cssStr)
}

func jQueryHandler(w http.ResponseWriter, r *http.Request) {
	jq, err := Asset("assets/jquery-2.2.0.min.js")
	if err != nil {
		log.Printf("Error %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/javascript")
	jqStr := string(jq)
	fmt.Fprint(w, jqStr)
}

func frm2jsHandler(w http.ResponseWriter, r *http.Request) {
	j, err := Asset("assets/form2js.js")
	if err != nil {
		log.Printf("Error %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/javascript")
	jStr := string(j)
	fmt.Fprint(w, jStr)
}

// apiResourceHandler is used by the dev-console ui to upate the resource map, triggering
// a websocket message to be sent to the bot process.
func apiResourceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}

	var postMsg Resource
	err = json.Unmarshal(body, &postMsg)
	if err != nil {
		log.Panicln("err: ", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "POST":
		registry.Add(postMsg.ID, &postMsg)

		// respond to http call
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(""))

	case "PATCH":
		updateResMsg := NewUpdateResourceMessage(postMsg.ID, &postMsg)
		registry.Update(updateResMsg)

		if wsConn != nil {
			wsConn.WriteJSON(updateResMsg)
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))

	case "DELETE":
		removeResMsg := NewRemoveResourceMessage(postMsg.ID)
		registry.Remove(removeResMsg)

		if wsConn != nil {
			wsConn.WriteJSON(removeResMsg)
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
	}

	log.Printf("registry:\n%s", registry)
}
