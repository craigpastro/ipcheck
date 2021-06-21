package storage

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var (
	DbPool     *pgxpool.Pool
	allMatches bool
	ipSetsDir  string
	ipSets     []string
)

type DbConfig struct {
	DatabaseURL string
	AllMatches  bool
	IPSetsDir   string
	IPSets      []string
}

type BlockedIP struct {
	Address    net.IP
	Blocklists []Blocklist
}

type Blocklist struct {
	Filename       string
	SourceFileDate string
}

func InitDb(config DbConfig) error {
	allMatches = config.AllMatches
	ipSetsDir = config.IPSetsDir
	ipSets = config.IPSets

	ctx := context.Background()
	var err error

	if DbPool, err = pgxpool.Connect(ctx, config.DatabaseURL); err != nil {
		return errors.Wrap(err, "unable to connect to the database")
	}

	schemaBytes, err := os.ReadFile("sql/schema.sql")
	if err != nil {
		return errors.Wrap(err, "error reading `sql/schema.sql`")
	}

	if _, err := DbPool.Exec(ctx, string(schemaBytes)); err != nil {
		return errors.Wrap(err, "failed to create the database schema")
	}

	return nil
}

func IsIPAddressInBlocklist(ipAddress net.IP) (*BlockedIP, error) {
	statement := "SELECT filename, source_file_date FROM blocklist WHERE address >>= $1"
	if !allMatches {
		statement = statement + " LIMIT 1"
	}

	rows, err := DbPool.Query(context.Background(), statement, ipAddress.String())
	if err != nil {
		return nil, errors.Wrap(err, "postgres error running query")
	}
	defer rows.Close()

	var blocklists []Blocklist
	var filename string
	var sourceFileDate string
	for rows.Next() {
		if err := rows.Scan(&filename, &sourceFileDate); err != nil {
			return nil, errors.Wrap(err, "postgres error scanning row")
		}
		blocklists = append(blocklists, Blocklist{filename, sourceFileDate})
	}

	if len(blocklists) > 0 {
		return &BlockedIP{ipAddress, blocklists}, nil
	}

	return nil, nil
}

func createTempTable() error {
	ctx := context.Background()

	if _, err := DbPool.Exec(ctx, "DROP TABLE IF EXISTS temp"); err != nil {
		return errors.Wrap(err, "postgres error dropping temp table")
	}

	if _, err := DbPool.Exec(ctx, `CREATE TABLE temp (address INET NOT NULL, filename TEXT NOT NULL, source_file_date TEXT)`); err != nil {
		return errors.Wrap(err, "postgres error creating temp table")
	}

	return nil
}

func addIPSetsToTempTable() error {
	ctx := context.Background()

	for _, ipSet := range ipSets {
		filename := filepath.Join(ipSetsDir, ipSet)
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("error reading ipset '%v': %v\n", filename, err)
		}
		defer file.Close()

		log.Printf("inserting '%v' into the temp table\n", filename)

		var sourceFileDate string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			l := scanner.Text()

			if !strings.HasPrefix(l, "#") {
				_, err = DbPool.Exec(ctx, `INSERT INTO temp VALUES ($1, $2, $3)`, l, ipSet, sourceFileDate)
				if err != nil {
					log.Printf("error inserting a row into the temp table: %v", err)
				}
			} else {
				if strings.HasPrefix(l, "# Source File Date") {
					sourceFileDate = strings.TrimSpace(strings.SplitN(l, ":", 2)[1])
				}
			}
		}
	}

	if _, err := DbPool.Exec(ctx, "CREATE INDEX IF NOT EXISTS address_idx ON temp USING GIST (address inet_ops)"); err != nil {
		log.Printf("error creating GiST index on temp table: %v\n", err)
	}

	return nil
}

func replaceBlocklistTableWithTempTable() error {
	ctx := context.Background()

	tx, err := DbPool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "postgres error creating the table replace transaction")
	}

	if _, err = tx.Exec(ctx, "DROP TABLE IF EXISTS blocklist"); err != nil {
		return errors.Wrap(err, "postgres error dropping the blocklist table")
	}

	if _, err = tx.Exec(ctx, "ALTER TABLE temp RENAME TO blocklist"); err != nil {
		return errors.Wrap(err, "postgres error renaming the temp table to blocklist")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "postgres error committing the table replace transaction")
	}

	return nil
}
