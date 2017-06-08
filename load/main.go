package main

import (
	"log"
	"time"

	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
	"github.com/peterhellberg/emojilib"
)

func main() {
	// Connect to Cassandra to make keyspace
	cluster := gocql.NewCluster("cassandra")
	cluster.Timeout = 10 * time.Second
	s, err := cluster.CreateSession()
	qe := gocassa.GoCQLSessionToQueryExecutor(s)
	c := gocassa.NewConnection(qe)
	err = c.DropKeySpace("emoji")
	if err != nil {
		log.Fatal(err)
	}
	err = c.CreateKeySpace("emoji")
	if err != nil {
		log.Fatal(err)
	}
	ks := c.KeySpace("emoji")

	type Search struct {
		Search, Name, Char string
	}
	tbl := ks.Table("emoji", &Search{}, gocassa.Keys{
		PartitionKeys:     []string{"Search"},
		ClusteringColumns: []string{"Name"},
	})
	tbl.CreateIfNotExist()

	for k, v := range emojilib.All() {
		var s Search

		// Insert direct searches: "smiley -> :), smiley"
		s.Search = k
		s.Char = v.Char
		s.Name = k
		err = tbl.Set(s).Run()
		if err != nil {
			log.Fatal("set: ", err)
		}

		// Insert keyword searches
		for _, kw := range v.Keywords {
			s.Search = kw
			err = tbl.Set(s).Run()
			if err != nil {
				log.Fatal("set: ", err)
			}
		}
	}
}
