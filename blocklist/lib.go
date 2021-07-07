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

var ranger cidranger.Ranger

func CloneRepoAndPopulateTrie(config BlConfig) error {
	if err := cloneRepo(config); err != nil {
		return err
	}

	PopulateTrie(config)
	return nil
}

func cloneRepo(config BlConfig) error {
	if err := os.RemoveAll(config.IPSetsDir); err != nil {
		return fmt.Errorf("error removing '%v': %w", config.IPSetsDir, err)
	}

	if _, err := git.PlainClone(config.IPSetsDir, false, &git.CloneOptions{URL: blocklistRepoURL}); err != nil {
		return fmt.Errorf("error cloning '%v': %w", blocklistRepoURL, err)
	}

	log.Printf("successfully cloned '%v'\n", blocklistRepoURL)
	return nil
}

func PopulateTrie(config BlConfig) {
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
					// Super hacky
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
