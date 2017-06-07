package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
	"github.com/peterhellberg/emojilib"
)

var eToName map[string]string

func init() {
	eToName = make(map[string]string, 1000)
	for k, v := range emojilib.All() {
		eToName[v.Char] = k
	}
}

// SearchResult is the top-level struct returned
// from a search in JSON format.
type SearchResult struct {
	Emojis []Emoji
	Error  string
}

type Emoji struct {
	Char string
	Name string
}

func SearchCass(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/search/") {
		http.Error(w, "no search", http.StatusBadRequest)
		return
	}

	result := &SearchResult{}
	in := r.URL.Path[8:]
	log.Println(in)

	type Search struct {
		Search, Name, Char string
	}
	tbl := ks.Table("emoji", &Search{}, gocassa.Keys{
		PartitionKeys:     []string{"Search"},
		ClusteringColumns: []string{"Name"},
	})

	res := make([]Search, 0, 100)
	if err := tbl.Where(gocassa.Eq("Search", in)).Read(&res).Run(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result.Emojis = make([]Emoji, len(res))
	for i, s := range res {
		result.Emojis[i].Char = s.Char
		result.Emojis[i].Name = s.Name
	}

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(out)
}

func SearchInternal(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/search/") {
		http.Error(w, "no search", http.StatusBadRequest)
		return
	}

	result := &SearchResult{}
	in := r.URL.Path[8:]

	// Try Find first. If we find it, return the exact
	// match.
	e, err := emojilib.Find(in)
	if err == nil {
		result.Emojis = []Emoji{{Char: e.Char, Name: in}}
	} else {
		// Otherwise try it as a keyword lookup.
		all, err := emojilib.Keyword(in)

		if err == nil {
			result.Emojis = make([]Emoji, len(all))
			for i, e2 := range all {
				result.Emojis[i].Char = e2.Char
				result.Emojis[i].Name = eToName[e2.Char]
			}
		} else {
			result.Error = err.Error()
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(out)
}

var ks gocassa.KeySpace

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	s, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	qe := gocassa.GoCQLSessionToQueryExecutor(s)
	ks = gocassa.NewConnection(qe).KeySpace("emoji")

	http.HandleFunc("/search/", SearchCass)

	addr := ":8000"
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
