package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Console struct {
	JQuery template.JS
	Css    template.CSS
	Data   string
	Logo   string
	Port   string
}

func NewConsole() (*Console, error) {
	// var data Data
	// strData, err := data.getFileData()
	// if err != nil {
	// 	return nil, err
	// }

	p := &Console{
		Port: Port,
		// Data: strData,
	}

	l, err := Asset("assets/beepboop_2x.ca5eff.png")
	if err != nil {
		log.Printf("Error %s", err.Error())
	}
	p.Logo = base64.StdEncoding.EncodeToString(l)

	return p, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p, err := NewConsole()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, p)
}

var indexTmpl, _ = Asset("assets/index.html")
var t = template.Must(template.New("index.html").Parse(string(indexTmpl)))

func renderTemplate(w http.ResponseWriter, p *Console) {
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
