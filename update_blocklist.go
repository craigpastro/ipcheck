package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

const blocklistRepoURL = "https://github.com/firehol/blocklist-ipsets"

func updateBlocklist() error {
	if err := cloneBlocklistRepo(); err != nil {
		log.Printf("error cloning blocklist repo: %v\n", err)
		return err
	}

	if err := createTempTable(); err != nil {
		log.Printf("error creating temp table: %v\n", err)
		return err
	}

	if err := addIPSetsToTempTable(); err != nil {
		log.Printf("error adding ipsets to temp table: %v\n", err)
		return err
	}

	if err := replaceBlocklistTableWithTempTable(); err != nil {
		log.Printf("error replacing blocklist table with temp table: %v\n", err)
		return err
	}

	return nil
}

func cloneBlocklistRepo() error {
	if err := os.RemoveAll(ipSetsDir); err != nil {
		return err
	}

	if _, err := git.PlainClone(ipSetsDir, false, &git.CloneOptions{URL: blocklistRepoURL}); err != nil {
		return err
	}

	return nil
}
