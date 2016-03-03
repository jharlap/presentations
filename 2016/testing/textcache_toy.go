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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/groupcache"
	"golang.org/x/net/context"
)

var (
	errFileNotFound = errors.New("file not found")
	textGroup       *groupcache.Group
	s3Srv           *s3.S3
)

const cacheLimitBytes = 10 << 20 // 1MB

func main() {
	me := os.Getenv("PEERS_ME")
	log.SetPrefix("[" + me + "] ")

	others := strings.Split(os.Getenv("PEERS_OTHER"), ",")
	all := append([]string{me}, others...)
	peers := groupcache.NewHTTPPool(all[0])
	peers.Set(all...)

	textGroup = groupcache.NewGroup("text", cacheLimitBytes, groupcache.GetterFunc(textGetter))
	s3Srv = s3.New(nil)

	http.HandleFunc("/reader/", handleGetText)

	u, err := url.Parse(me)
	if err != nil {
		log.Fatalf("Error parsing PEERS_ME: %s", err)
	}
	p := strings.Split(u.Host, ":")[1]
	log.Println("Listening on port", p)
	log.Println(http.ListenAndServe("0.0.0.0:"+p, nil))
}

func handleGetText(rw http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/reader/")
	var data groupcache.ByteView
	err := textGroup.Get(context.TODO(), key, groupcache.ByteViewSink(&data))
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

func textGetter(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	filename := parseKey(key)
	data, err := fetch(filename)
	if err != nil {
		return err
	}
	return dest.SetBytes(data)
}

func parseKey(k string) string {
	if ks := strings.SplitN(k, "-", 2); len(ks) > 0 {
		return ks[0]
	}
	return k
}

func fetch(f string) ([]byte, error) {
	log.Println("Fetching text", f)
	if len(f) < 4 {
		return nil, errFileNotFound
	}

	key := "text/" + f[:len(f)-3] + "/" + f + "_par_id.txt"
	r, err := s3Srv.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("wattpad.staging"),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}
