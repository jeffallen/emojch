package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var dev = flag.Bool("dev", false, "Run in dev mode.")

var result = template.Must(template.New("result").Parse(string(MustAsset("static/result.tmpl"))))

type SearchResult struct {
	Emojis []Emoji
	Error  string
}

type Emoji struct {
	Char string
	Name string
}

func Search(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("http://search/%s", url.PathEscape(r.FormValue("q")))
	log.Println("search", url)
	resp, err := http.Get(url)
	var all []byte
	if err == nil {
		all, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	var sr SearchResult
	if err != nil {
		log.Println("bad query to search: ", err)
	} else {
		err = json.Unmarshal(all, &sr)
		if err != nil {
			log.Println("bad search result: ", err)
		}
	}

	result.Execute(w, sr)
}

func main() {
	flag.Parse()

	if *dev {
		http.Handle("/", http.FileServer(http.Dir("static")))
	} else {
		http.Handle("/", http.FileServer(assetFS()))
	}

	http.HandleFunc("/search", Search)

	addr := ":8000"
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
