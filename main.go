package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

var addr = flag.String("addr", ":8999", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

var validPath = regexp.MustCompile("^/ws/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		m := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Println(r)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		serveWs(w, r, m[1])
	}
}

func main() {
	flag.Parse()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws/", makeHandler(serveWs))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
