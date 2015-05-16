package testutil

import(
	"fmt"
	"regexp"
	"net/http"
	"net/http/httptest"
)

type Mockconsul struct {
	Host string
	Port string
	ts *httptest.Server
} 

func (m *Mockconsul ) MockConsul() error {
	// Test server that always hands back jsondata
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jsondata string
		if r.URL.Path == "/v1/catalog/nodes" {
			w.WriteHeader(200)
			jsondata = "[{\"Node\":\"baz\",\"Address\":\"10.1.10.11\"},{\"Node\":\"foobar\",\"Address\":\"10.1.10.12\"}]" 
		
		} else if r.URL.Path == "/v1/catalog/node/baz" {
			w.WriteHeader(200)
			jsondata = "{\"bogus\" :{\"ID\":\"redis\",\"Service\":\"redis\",\"Tags\": [],\"Port\": 5641}}"
		} else if r.URL.Path == "/v1/catalog/node/foobar" {
			w.WriteHeader(200)
			jsondata = "{}"
		} else if r.URL.Path == "/v1/catalog/services" {
			w.WriteHeader(200)
			jsondata = "{\"redis\":[]}"
		} else if r.URL.Path == "/v1/catalog/service/redis" {
			w.WriteHeader(200)
			jsondata = "[{\"Node\":\"bogus\",\"Address\":\"127.0.0.1\",\"ServiceID\":\"redis\",\"ServiceName\":\"redis\",\"ServiceTags\":[],\"ServiceAddress\":\"\",\"ServicePort\":5641}]"
		} else {
			w.WriteHeader(500)
			jsondata = "[{}]"
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w,jsondata)
	}))
	m.ts = ts
	// Get the address and port of the test server
	re,err := regexp.Compile(`http://([^:]+):(\d+).*`)
	if err != nil {
		return err
	}
	result := re.FindStringSubmatch(ts.URL)
	m.Host = result[1]
	m.Port = result[2]
	return nil
}

func (m *Mockconsul) Close() {
	m.ts.Close()
}