package main

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	dbConfig := dbConfig{
		"postgres://postgres:password@127.0.0.1:5432/test",
		true,
		"test_ipsets",
		[]string{"a.netset", "b.ipset", "c.ipset"},
	}

	initDb(dbConfig)
	if err := updateBlocklists(); err != nil {
		checkError(err, "error updating blocklist")
	}

	os.Exit(m.Run())
}

func TestIPAddressIsNotInABlocklist(t *testing.T) {
	ipAddress := "1.2.3.4"

	r := setupRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestIPAddressIsInANetSet(t *testing.T) {
	ipAddress := "5.188.206.37"
	want := Response{
		net.ParseIP(ipAddress),
		[]Blocklist{{"a.netset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := setupRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp Response
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestIPAddressIsInAnIPSet(t *testing.T) {
	ipAddress := "1.0.175.32"
	want := Response{
		net.ParseIP(ipAddress),
		[]Blocklist{{"b.ipset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := setupRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp Response
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestIPAddressIsInTwoBlocklists(t *testing.T) {
	ipAddress := "45.134.26.37"
	want := Response{
		net.ParseIP(ipAddress),
		[]Blocklist{
			{"a.netset", "Thu Jun 17 20:30:46 UTC 2021"},
			{"b.ipset", "Thu Jun 17 20:30:46 UTC 2021"},
		},
	}

	r := setupRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp Response
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestIPV6AddressIsNotInABlocklist(t *testing.T) {
	ipAddress := "::1"

	r := setupRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestIPV6AddressIsInAnIPSet(t *testing.T) {
	ipAddress := "2001:db8::8a2e:370:7334"
	want := Response{
		net.ParseIP(ipAddress),
		[]Blocklist{{"c.ipset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := setupRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp Response
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}

func TestIPV6AddressIsInANetSet(t *testing.T) {
	ipAddress := "2021:abc:d::1"
	want := Response{
		net.ParseIP(ipAddress),
		[]Blocklist{{"c.ipset", "Thu Jun 17 20:30:46 UTC 2021"}},
	}

	r := setupRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	var resp Response
	json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, resp)
}
