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
	"github.com/pkg/errors"
	"github.com/yl2chen/cidranger"
)

const blocklistRepoURL = "https://github.com/firehol/blocklist-ipsets"

type BlConfig struct {
	IPSetsDir string
	IPSets    []string
}

var ranger cidranger.Ranger

func InitBl(config BlConfig, cloneRepo bool) {
	if cloneRepo {
		cloneBlocklistRepo(config)
	}

	ranger = cidranger.NewPCTrieRanger()

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
				_, network, err := net.ParseCIDR(l)
				if err != nil {
					// Super hacky
					_, network, err = net.ParseCIDR(l + "/32")
				}

				if err := ranger.Insert(cidranger.NewBasicRangerEntry(*network)); err != nil {
					log.Printf("error inserting '%v' in the trie: %v", network, err)
					continue
				}
				log.Printf("inserted '%v' in the trie", network)
			}
		}
	}
}

func cloneBlocklistRepo(config BlConfig) error {
	if err := os.RemoveAll(config.IPSetsDir); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error removing '%v'", config.IPSetsDir))
	}

	if _, err := git.PlainClone(config.IPSetsDir, false, &git.CloneOptions{URL: blocklistRepoURL}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error cloning '%v'", blocklistRepoURL))
	}

	log.Printf("successfully cloned '%v'\n", blocklistRepoURL)

	return nil
}

func InBlocklist(ip net.IP) (bool, error) {
	res, err := ranger.Contains(ip)
	if err != nil {
		return false, errors.Wrap(err, "error checking containment in the trie")
	}

	return res, nil
}
