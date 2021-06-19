package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

const blocklistRepoURL = "https://github.com/firehol/blocklist-ipsets"

func updateBlocklist() error {
	if err := cloneBlocklistRepo(); err != nil {
		return errors.Wrap(err, "error cloning blocklist repo")
	}

	if err := createTempTable(); err != nil {
		return errors.Wrap(err, "error creating temp table")
	}

	if err := addIPSetsToTempTable(); err != nil {
		return errors.Wrap(err, "error adding ipsets to temp table")
	}

	if err := replaceBlocklistTableWithTempTable(); err != nil {
		return errors.Wrap(err, "error replacing blocklist table with temp table")
	}

	log.Println("successfully updated blocklist table")
	return nil
}

func cloneBlocklistRepo() error {
	if err := os.RemoveAll(ipSetsDir); err != nil {
		return errors.Wrap(err, "error removing current IP_SETS_DIR")
	}

	if _, err := git.PlainClone(ipSetsDir, false, &git.CloneOptions{URL: blocklistRepoURL}); err != nil {
		return errors.Wrap(err, "error cloning "+blocklistRepoURL)
	}

	return nil
}
