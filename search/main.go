package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

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

func Search(w http.ResponseWriter, r *http.Request) {
	result := &SearchResult{}
	in := r.URL.Path
	// Remove initial slash.
	if len(in) > 1 && in[0] == '/' {
		in = in[1:]
	}
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
		result.Error = err.Error()
	} else {
		result.Emojis = make([]Emoji, len(res))
		for i, s := range res {
			result.Emojis[i].Char = s.Char
			result.Emojis[i].Name = s.Name
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
	cluster := gocql.NewCluster("cassandra")
	s, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	qe := gocassa.GoCQLSessionToQueryExecutor(s)
	ks = gocassa.NewConnection(qe).KeySpace("emoji")

	http.HandleFunc("/", Search)

	addr := ":8000"
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
