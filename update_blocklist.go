package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

func updateBlocklist() error {
	err := cloneOrPullBlockListRepo()
	if err != nil {
		log.Printf("clone or pull blocklist repo: %v\n", err)
	}

	failures := addIPSetsToTempTable()
	if failures > 0 {
		log.Printf("error adding ipsets to temp table: %v\n", err)
	}

	replaceBlocklistTableWithTempTable()

	return nil
}

func cloneOrPullBlockListRepo() error {
	var err error

	if _, err = os.Stat(ipSetsDir); os.IsNotExist(err) {
		log.Println("cloning blocklist repo")

		_, err = git.PlainClone(ipSetsDir, false, &git.CloneOptions{
			URL: "https://github.com/firehol/blocklist-ipsets",
		})
	} else {
		log.Println("pulling blocklist repo")

		r, err := git.PlainOpen(ipSetsDir)
		if err != nil {
			return err
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			return err
		}
	}

	return err
}
