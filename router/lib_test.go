package router

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/siyopao/ipcheck/storage"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	dbConfig := storage.DbConfig{
		DatabaseURL: "postgres://postgres:password@127.0.0.1:6543/test",
		AllMatches:  true,
		IpSetsDir:   "../test_ipsets",
		IpSets:      []string{"a.netset", "b.ipset", "c.ipset"},
	}

	storage.InitDb(dbConfig)
	if err := storage.UpdateBlocklists(); err != nil {
		fmt.Println("error updating the blocklists")
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestIPAddressIsNotInABlocklist(t *testing.T) {
	ipAddress := "1.2.3.4"

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestIPAddressIsInANetSet(t *testing.T) {
	ipAddress := "5.188.206.37"
	want := BlockedIP{
		net.ParseIP(ipAddress),
		[]Blocklist{{"a.netset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp BlockedIP
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestIPAddressIsInAnIPSet(t *testing.T) {
	ipAddress := "1.0.175.32"
	want := BlockedIP{
		net.ParseIP(ipAddress),
		[]Blocklist{{"b.ipset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp BlockedIP
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestIPAddressIsInTwoBlocklists(t *testing.T) {
	ipAddress := "45.134.26.37"
	want := BlockedIP{
		net.ParseIP(ipAddress),
		[]Blocklist{
			{"a.netset", "Thu Jun 17 20:30:46 UTC 2021"},
			{"b.ipset", "Thu Jun 17 20:30:46 UTC 2021"},
		},
	}

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp BlockedIP
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestNotAllMatchesReturnsOneMatch(t *testing.T) {
	// Need to resetup the database pool for this test.
	dbConfig := storage.DbConfig{
		DatabaseURL: "postgres://postgres:password@127.0.0.1:6543/test",
		AllMatches:  false,
		IpSetsDir:   "test_ipsets",
		IpSets:      []string{"a.netset", "b.ipset", "c.ipset"},
	}
	storage.InitDb(dbConfig)

	ipAddress := "45.134.26.37"

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp BlockedIP
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, ipAddress, resp.IPAddress.String())
	assert.Equal(t, 1, len(resp.Blocklists))
}

func TestIPV6AddressIsNotInABlocklist(t *testing.T) {
	ipAddress := "::1"

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestIPV6AddressIsInAnIPSet(t *testing.T) {
	ipAddress := "2001:db8::8a2e:370:7334"
	want := BlockedIP{
		net.ParseIP(ipAddress),
		[]Blocklist{{"c.ipset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp BlockedIP
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestIPV6AddressIsInANetSet(t *testing.T) {
	ipAddress := "2021:abc:d::1"
	want := BlockedIP{
		net.ParseIP(ipAddress),
		[]Blocklist{{"c.ipset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp BlockedIP
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestNotAnIPAddress(t *testing.T) {
	ipAddress := "foo"

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
