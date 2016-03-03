package main

import (
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"

	"testing/textcache/2/textcache"
)

var (
	originHost string
)

func init() {
	os.Setenv("PEERS_ME", "http://localhost:9001")
	os.Setenv("PEERS_OTHER", "http://localhost:9002,http://localhost:9003")
	os.Setenv("ORIGIN", "http://localhost:4545")
}

func main() {
	originHost = os.Getenv("ORIGIN")
	me := os.Getenv("PEERS_ME")
	others := strings.Split(os.Getenv("PEERS_OTHER"), ",")
	log.SetPrefix("[" + me + "] ")

	all := append([]string{me}, others...)

	cache := textcache.New(originHost, me) // HL
	cache.UpdatePeers(all)                 // HL
	http.Handle("/text/", cache)           // HL

	u, err := url.Parse(me)
	if err != nil {
		log.Fatalf("Error parsing PEERS_ME: %s", err)
	}
	p := strings.Split(u.Host, ":")[1]
	log.Println("Listening on port", p)
	log.Println(http.ListenAndServe("0.0.0.0:"+p, nil))
}
