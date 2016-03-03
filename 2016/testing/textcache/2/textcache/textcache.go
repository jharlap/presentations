package textcache

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang/groupcache"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const cacheLimitBytes = 10 << 20 // 1MB

var errFileNotFound = errors.New("file not found")

type Service struct {
	originHost string
	pool       *groupcache.HTTPPool
	group      *groupcache.Group
}

func New(originHost, myselfHost string) *Service {
	s := &Service{
		originHost: originHost,
		pool:       groupcache.NewHTTPPool(myselfHost),
	}
	s.group = groupcache.NewGroup("text", cacheLimitBytes, groupcache.GetterFunc(s.textGetter))
	return s
}

var testPool *groupcache.HTTPPool
const testMyselfHost = "http://me.example.com"

func NewTestService(originHost string) *Service {
	if testPool == nil {
		testPool = groupcache.NewHTTPPool(testMyselfHost)
	}
	s := &Service{
		originHost: originHost,
		pool:       testPool,
	}
	s.group = groupcache.NewGroup("text_"+originHost, cacheLimitBytes, groupcache.GetterFunc(s.textGetter))
	return s
}

func (s *Service) UpdatePeers(p []string) {
	s.pool.Set(p...)
}

func (s *Service) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var data groupcache.ByteView
	err := s.group.Get(context.TODO(), r.URL.Path, groupcache.ByteViewSink(&data)) // HL
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

func (s *Service) textGetter(gctx groupcache.Context, key string, dest groupcache.Sink) error {
	ctx := gctx.(context.Context)

	url := s.originHost + originPathForKey(key)
	if len(url) < 10 {
		return errFileNotFound
	}

	data, err := fetch(ctx, url) // HL
	if err != nil {
		return err
	}
	return dest.SetBytes(data) // HL
}

func originPathForKey(k string) string {
	ss := strings.Split(k, "/")
	return strings.Join(ss[:len(ss)-1], "/")
}

func fetch(ctx context.Context, url string) ([]byte, error) {
	log.Println("Fetching text", url)
	r, err := ctxhttp.Get(ctx, nil, url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}
