package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/siyopao/ipcheck/blocklist"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	blocklist.PopulateTrie(blocklist.BlConfig{
		IPSetsDir: "../test_ipsets",
		IPSets:    []string{"a.netset", "b.ipset", "c.ipset"},
	})

	os.Exit(m.Run())
}

func TestIPAddressIsNotInABlocklist(t *testing.T) {
	ipAddress := "1.2.3.4"
	want := Response{false}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)

}

func TestIPAddressIsInANetSet(t *testing.T) {
	ipAddress := "5.188.206.37"
	want := Response{true}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPAddressIsInAnIPSet(t *testing.T) {
	ipAddress := "1.0.175.32"
	want := Response{true}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPAddressIsInTwoBlocklists(t *testing.T) {
	ipAddress := "45.134.26.37"
	want := Response{true}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestNotAllMatchesReturnsOneMatch(t *testing.T) {
	ipAddress := "45.134.26.37"
	want := Response{true}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPV6AddressIsNotInABlocklist(t *testing.T) {
	ipAddress := "::1"
	want := Response{false}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPV6AddressIsInAnIPSet(t *testing.T) {
	ipAddress := "2001:db8::8a2e:370:7334"
	want := Response{true}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestIPV6AddressIsInANetSet(t *testing.T) {
	ipAddress := "2021:abc:d::1"
	want := Response{true}
	var got Response

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, want, got)
}

func TestNotAnIPAddress(t *testing.T) {
	ipAddress := "foo"

	r := InitRouter("release")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/addresses/"+ipAddress, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
