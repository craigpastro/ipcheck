package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/siyopao/ipcheck/storage"
)

const blocklistRepoURL = "https://github.com/firehol/blocklist-ipsets"

func updateBlocklists() error {
	if err := cloneBlocklistRepo(); err != nil {
		return errors.Wrap(err, "error cloning blocklist repo")
	}

	if err := storage.CreateTempTable(); err != nil {
		return errors.Wrap(err, "error creating temp table")
	}

	if err := storage.AddIPSetsToTempTable(); err != nil {
		return errors.Wrap(err, "error adding ipsets to temp table")
	}

	if err := storage.ReplaceBlocklistTableWithTempTable(); err != nil {
		return errors.Wrap(err, "error replacing blocklist table with temp table")
	}

	log.Println("successfully updated blocklist table")
	return nil
}

func cloneBlocklistRepo() error {
	if err := os.RemoveAll(storage.IpSetsDir); err != nil {
		return errors.Wrap(err, "error removing current IP_SETS_DIR")
	}

	if _, err := git.PlainClone(storage.IpSetsDir, false, &git.CloneOptions{URL: blocklistRepoURL}); err != nil {
		return errors.Wrap(err, "error cloning "+blocklistRepoURL)
	}

	return nil
}
