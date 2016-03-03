package textcache_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"testing/textcache/2/textcache"
)

func TestFetchFromOrigin(t *testing.T) {
	ts := httptest.NewServer(fakeOrigin(1))
	defer ts.Close()

	c := textcache.NewTestService(ts.URL)

	// request a text from origin 1
	rec := httptest.NewRecorder()
	req := fakeRequest(ts.URL + "/text/123456/abcdef/1")
	c.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("request code == %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "origin 1 at") {
		t.Errorf("request body == '%s', want contains 'origin 1 at'", body)
	}
}

func TestSecondRequestFromCache(t *testing.T) {
	ts := httptest.NewServer(fakeOrigin(1))
	defer ts.Close()

	c := textcache.NewTestService(ts.URL)

	// request a text from origin
	rec := httptest.NewRecorder()
	req := fakeRequest(ts.URL + "/text/123456/abcdef/1")
	c.ServeHTTP(rec, req)

	body1 := rec.Body.String()

	// 2nd request should be from the cache and thus contain the same timestamp as the 1st
	time.Sleep(5 * time.Millisecond)
	rec = httptest.NewRecorder()
	c.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("2 request code == %d, want 200", rec.Code)
	}

	body2 := rec.Body.String()
	if body1 != body2 {
		t.Errorf("2 body == '%s', want '%s'", body2, body1)
	}
}

func TestFakeOrigin(t *testing.T) {
	fo := fakeOrigin(1)

	// make a first request
	rec := httptest.NewRecorder()
	req := fakeRequest("http://example.com/text/123456/abcdef/1")
	fo.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("1 request code == %d, want 200", rec.Code)
	}

	body1 := rec.Body.String()
	if !strings.Contains(body1, "origin 1 at") {
		t.Errorf("1 request body == '%s', want contains 'origin 1 at'", body1)
	}

	// request the same resource again after a few ms to get a different timestamp in the body
	time.Sleep(5 * time.Millisecond)
	rec = httptest.NewRecorder()
	fo.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("2 request code == %d, want 200", rec.Code)
	}

	body2 := rec.Body.String()
	if !strings.Contains(body2, "origin 1 at") {
		t.Errorf("2 request body == '%s', want contains 'origin 1 at'", body2)
	}

	if body1 == body2 {
		t.Errorf("2 body timestamp should be different. got '%s'", body2)
	}
}

type fakeOrigin int

func (f fakeOrigin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "an amazing story\nserved by origin %d at time %d\n", f, time.Now().UnixNano())
}

func fakeRequest(URL string) *http.Request {
	u, _ := url.Parse(URL)
	return &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
}
