package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {

	// Needs at least 2 args including the ./main
	if len(os.Args) < 2 {
		fmt.Println("usage: ./main <insert|select|select-bypass-cache> [ addr1, addr2, ...]")
		return
	}

	action := os.Args[1]
	addrList := []string{"127.0.0.1:9042"}

	// If the addresses are passed, use them instead
	if len(os.Args) >= 3 {
		addrList = os.Args[2:]
	}

	cluster := gocql.NewCluster(addrList...)
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}

	// Connect to cassandra/scylla
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	start := time.Now()
	if action == "insert" {
		insert_data(session)
	} else if action == "select-bypass-cache" {
		select_data(session, true)
	} else {
		select_data(session, false)
	}
	log.Printf("Elapsed time %s", time.Since(start))

}

func insert_data(session *gocql.Session) {

	// Create some context
	ctx := context.Background()

	// Drop the keyspace if it already exists
	err := session.Query(`DROP KEYSPACE IF EXISTS demo`).WithContext(ctx).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new keyspace
	err = session.Query(
		`CREATE KEYSPACE demo WITH REPLICATION = { 
      'class': 'NetworkTopologyStrategy',
      'replication_factor': 3
    }`,
	).WithContext(ctx).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Create a table
	err = session.Query(
		`CREATE TABLE demo.flight (
      uniquecarrier text,
      flightnum text,
      tailnum text,
      actualelapsedtime text,
      crselapsedtime text,
      distance text,
      PRIMARY KEY(uniquecarrier, flightnum, distance)
    )`,
	).WithContext(ctx).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Open the file
	file, err := os.Open("./airline.csv.sample")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create batch and template statement
	batch := session.NewBatch(gocql.LoggedBatch)
	stmt := `INSERT INTO demo.flight (uniquecarrier, flightnum, tailnum, actualelapsedtime, crselapsedtime, distance) VALUES (?, ?, ?, ?, ?, ?);`

	// Create a reader and scanner with ISO 8859-1 encoding
	reader := bufio.NewReader(
		transform.NewReader(
			file,
			charmap.ISO8859_1.NewDecoder(),
		),
	)
	scanner := bufio.NewScanner(reader)

	// Define batch vars and counters
	count := 0
	batchCount := 0
	batchSize := 10

	// For each line of the file
	for scanner.Scan() {

		line := strings.Split(scanner.Text(), ",")
		actualelapsedtime := line[0]
		crselapsedtime := line[6]
		distance := line[15]
		flightnum := line[17]
		tailnum := line[23]
		uniquecarrier := line[26]

		// Craft a batch query with the template statement
		batch.Query(
			stmt,
			uniquecarrier,
			flightnum,
			tailnum,
			actualelapsedtime,
			crselapsedtime,
			distance,
		)
		count++

		// Execute batch when it reaches the batch size or at the end of the file
		if count%batchSize == 0 {
			err = session.ExecuteBatch(batch)
			if err != nil {
				log.Fatal(err)
			}
			batchCount++
			fmt.Printf("inserted batchNo: %d with %d records\n", batchCount, count)
			count = 0
			batch = session.NewBatch(gocql.LoggedBatch)
		}
	}

	// Clean up the last batch
	err = session.ExecuteBatch(batch)
	if err != nil {
		log.Panic(err)
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func select_data(session *gocql.Session, bypassCache bool) {
	// Execute the SELECT query
	query := "SELECT * FROM demo.flight WHERE uniquecarrier='MQ'"
	if bypassCache {
		query += " BYPASS CACHE;"
	}
	iter := session.Query(query).Iter()

	// Iterate over the results
	var uniquecarrier, flightnum, tailnum, actualelapsedtime, crselapsedtime, distance string

	counter := 0
	// For each line, print it out
	for iter.Scan(
		&uniquecarrier,
		&flightnum,
		&tailnum,
		&actualelapsedtime,
		&crselapsedtime,
		&distance,
	) {
		// Process each row
		fmt.Printf("%s, %s, %s, %s, %s, %s\n",
			uniquecarrier,
			flightnum,
			tailnum,
			actualelapsedtime,
			crselapsedtime,
			distance,
		)
		counter++
	}

	// Check for errors
	if err := iter.Close(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Read total of %d records:", counter)
}
