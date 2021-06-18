package main

import (
	// "bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

func initDb() {
	ctx := context.Background()

	var err error
	dbPool, err = pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	checkError("unable to connect to the database", err)

	schemaBytes, err := os.ReadFile("sql/schema.sql")
	checkError("error reading `sql/schema.sql`", err)

	_, err = dbPool.Exec(ctx, string(schemaBytes))
	checkError("failed to create the database schema", err)
}

func checkBlocklist(ipAddress net.IP) bool {
	statement := `SELECT address FROM blocklist WHERE address >>= $1`

	var res net.IP
	err := dbPool.QueryRow(context.Background(), statement, ipAddress.String()).Scan(&res)

	if err == nil {
		return true
	}

	return false
}

func addIPSetsToTempTable() int {
	ipSets := strings.Split(os.Getenv("IP_SETS"), ",")
	// var filesCopied int

	for _, fileName := range ipSets {
		file, err := os.Open(fileName)
		if err != nil {
			log.Printf("error reading ipset `%v`\n", fileName)
		}
		defer file.Close()

		// TODO: Get this working.
		// reader := bufio.NewReader(file)
		// _, err = dbPool.CopyFrom(context.Background(), pgx.Identifier{"temp"}, []string{"address"}, reader, pgx.CopyFromRows(reader))
		// if err == nil {
		// 	filesCopied++
		// }
		// checkError(">>>>>", err)
	}

	return 0 // len(ipSets) - filesCopied // The number of 'copy's that failed.
}

func replaceBlocklistTableWithTempTable() error {
	// TODO: Implement this.

	return nil
}
