package blocklist

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/yl2chen/cidranger"
)

const blocklistRepoURL = "https://github.com/firehol/blocklist-ipsets"

type BlConfig struct {
	IPSetsDir string
	IPSets    []string
}

var Config BlConfig
var ranger cidranger.Ranger

func InitBlocklists(config BlConfig) error {
	Config = config

	if err := cloneRepo(Config.IPSetsDir); err != nil {
		return err
	}

	populateTrie(Config)
	return nil
}

func cloneRepo(ipSetsDir string) error {
	if err := os.RemoveAll(ipSetsDir); err != nil {
		return fmt.Errorf("error removing '%v': %w", ipSetsDir, err)
	}

	if _, err := git.PlainClone(ipSetsDir, false, &git.CloneOptions{URL: blocklistRepoURL}); err != nil {
		return fmt.Errorf("error cloning '%v': %w", blocklistRepoURL, err)
	}

	log.Printf("successfully cloned '%v'\n", blocklistRepoURL)
	return nil
}

func populateTrie(config BlConfig) {
	newRanger := cidranger.NewPCTrieRanger()

	for _, ipSet := range config.IPSets {
		filename := filepath.Join(config.IPSetsDir, ipSet)
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("error reading ipset '%v': %v\n", filename, err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			l := scanner.Text()

			if !strings.HasPrefix(l, "#") {
				ip, network, err := net.ParseCIDR(l)
				if err != nil {
					ip = net.ParseIP(l)
					if ip == nil {
						continue
					}

					if ip.To4() != nil {
						// So an IPV4 address.
						_, network, _ = net.ParseCIDR(l + "/32")
					} else {
						_, network, _ = net.ParseCIDR(l + "/128")
					}
				}

				if err := newRanger.Insert(cidranger.NewBasicRangerEntry(*network)); err != nil {
					log.Printf("error inserting '%v' in the trie: %v", network, err)
					continue
				}
			}
		}
	}

	ranger = newRanger
	log.Printf("trie has been populated")
}

func InBlocklist(ip net.IP) (bool, error) {
	res, err := ranger.Contains(ip)
	if err != nil {
		return false, fmt.Errorf("error checking containment in the trie: %w", err)
	}

	return res, nil
}
