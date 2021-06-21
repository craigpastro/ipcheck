package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

const blocklistRepoURL = "https://github.com/firehol/blocklist-ipsets"

func CloneAndUpdateBlocklists() error {
	if err := cloneBlocklistRepo(); err != nil {
		return errors.Wrap(err, "error cloning blocklist repo")
	}

	return UpdateBlocklists()
}

func cloneBlocklistRepo() error {
	if err := os.RemoveAll(ipSetsDir); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error removing '%v'", ipSetsDir))
	}

	if _, err := git.PlainClone(ipSetsDir, false, &git.CloneOptions{URL: blocklistRepoURL}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error cloning '%v'", blocklistRepoURL))
	}

	log.Printf("successfully cloned '%v'\n", blocklistRepoURL)

	return nil
}

// Visible for testing.
func UpdateBlocklists() error {
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
