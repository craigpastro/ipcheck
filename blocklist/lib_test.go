package blocklist

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config BlConfig

func TestMain(m *testing.M) {
	config = BlConfig{
		IPSetsDir: "../test_ipsets",
		IPSets:    []string{"a.netset", "b.ipset", "c.ipset"},
	}
	populateTrie(config)

	os.Exit(m.Run())
}

func TestIPAddressIsNotInABlocklist(t *testing.T) {
	resp, _ := InBlocklist(net.ParseIP("1.2.3.4"))
	assert.False(t, resp)
}

func TestIPAddressIsInANetSet(t *testing.T) {
	resp, _ := InBlocklist(net.ParseIP("5.188.206.37"))
	assert.True(t, resp)
}

func TestIPAddressIsInAnIPSet(t *testing.T) {
	resp, _ := InBlocklist(net.ParseIP("1.0.175.32"))
	assert.True(t, resp)
}

func TestIPAddressInTwoBlocklists(t *testing.T) {
	resp, _ := InBlocklist(net.ParseIP("45.134.26.37"))
	assert.True(t, resp)
}

func TestIPV6AddressIsNotInABlocklist(t *testing.T) {
	resp, _ := InBlocklist(net.ParseIP("::1"))
	assert.False(t, resp)
}

func TestIPV6AddressIsInAnIPSet(t *testing.T) {
	resp, _ := InBlocklist(net.ParseIP("2001:db8::8a2e:370:7334"))
	assert.True(t, resp)
}

func TestIPV6AddressIsInANetSet(t *testing.T) {
	resp, _ := InBlocklist(net.ParseIP("2021:abc:d::1"))
	assert.True(t, resp)
}
