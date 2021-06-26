package router

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/siyopao/ipcheck/blocklist"
	"github.com/stretchr/testify/assert"
)

var config blocklist.BlConfig

func TestMain(m *testing.M) {
	config = blocklist.BlConfig{
		IPSetsDir: "../test_ipsets",
		IPSets:    []string{"a.netset", "b.ipset", "c.ipset"},
	}
	blocklist.PopulateTrie(config)

	os.Exit(m.Run())
}

func TestIPAddressIsNotInABlocklist(t *testing.T) {
	ip := "1.2.3.4"
	want := Response{net.ParseIP(ip), false}
	var got Response

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r := InitRouter("release", config)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)

}

func TestIPAddressIsInANetSet(t *testing.T) {
	ip := "5.188.206.37"
	want := Response{net.ParseIP(ip), true}
	var got Response

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPAddressIsInAnIPSet(t *testing.T) {
	ip := "1.0.175.32"
	want := Response{net.ParseIP(ip), true}
	var got Response

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPAddressIsInTwoBlocklists(t *testing.T) {
	ip := "45.134.26.37"
	want := Response{net.ParseIP(ip), true}
	var got Response

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestNotAllMatchesReturnsOneMatch(t *testing.T) {
	ip := "45.134.26.37"
	want := Response{net.ParseIP(ip), true}
	var got Response

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPV6AddressIsNotInABlocklist(t *testing.T) {
	ip := "::1"
	want := Response{net.ParseIP(ip), false}
	var got Response

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPV6AddressIsInAnIPSet(t *testing.T) {
	ip := "2001:db8::8a2e:370:7334"
	want := Response{net.ParseIP(ip), true}
	var got Response

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPV6AddressIsInANetSet(t *testing.T) {
	ip := "2021:abc:d::1"
	want := Response{net.ParseIP(ip), true}
	var got Response

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestNotAnIPAddress(t *testing.T) {
	ip := "foo"

	r := InitRouter("release", config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ip, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
