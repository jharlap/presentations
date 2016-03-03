package main

import (
	"errors"
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"io/ioutil"
	"net/url"

	"github.com/golang/groupcache"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

var (
	errFileNotFound = errors.New("file not found")
	textGroup       *groupcache.Group
	originHost      string
)

const cacheLimitBytes = 10 << 20 // 1MB

func init() {
	os.Setenv("PEERS_ME", "http://localhost:9001")
	os.Setenv("PEERS_OTHER", "http://localhost:9002,http://localhost:9003")
	os.Setenv("ORIGIN", "http://localhost:4545")
}

func main() {
	originHost = os.Getenv("ORIGIN")
	me := os.Getenv("PEERS_ME")
	log.SetPrefix("[" + me + "] ")

	others := strings.Split(os.Getenv("PEERS_OTHER"), ",")
	all := append([]string{me}, others...)
	peers := groupcache.NewHTTPPool(all[0])
	peers.Set(all...)

	textGroup = groupcache.NewGroup("text", cacheLimitBytes, // HL
		groupcache.GetterFunc(textGetter)) // HL
	http.HandleFunc("/text/", handleGetText) // HL

	u, err := url.Parse(me)
	if err != nil {
		log.Fatalf("Error parsing PEERS_ME: %s", err)
	}
	p := strings.Split(u.Host, ":")[1]
	log.Println("Listening on port", p)
	log.Println(http.ListenAndServe("0.0.0.0:"+p, nil))
}

func handleGetText(rw http.ResponseWriter, r *http.Request) {
	var data groupcache.ByteView
	err := textGroup.Get(context.TODO(), r.URL.Path, groupcache.ByteViewSink(&data)) // HL
	if err == errFileNotFound {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeContent(rw, r, "", time.Now(), data.Reader())
}

func textGetter(gctx groupcache.Context, key string, dest groupcache.Sink) error {
	ctx := gctx.(context.Context)
	data, err := fetch(ctx, originPathForKey(key)) // HL
	if err != nil {
		return err
	}
	return dest.SetBytes(data) // HL
}

func originPathForKey(k string) string {
	ss := strings.Split(k, "/")
	return strings.Join(ss[:len(ss)-1], "/")
}

func fetch(ctx context.Context, originPath string) ([]byte, error) {
	log.Println("Fetching text", originPath)
	if len(originPath) < 10 {
		return nil, errFileNotFound
	}

	r, err := ctxhttp.Get(ctx, nil, originHost+originPath)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}
