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

	cluster := gocql.NewCluster("127.0.0.1:9042")
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

	// Create some context
	ctx := context.Background()

	// Drop the keyspace if it already exists
	err = session.Query(`DROP KEYSPACE IF EXISTS demo`).WithContext(ctx).Exec()
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
      PRIMARY KEY(uniquecarrier, flightnum, tailnum)
    )`,
	).WithContext(ctx).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Open the file
	start := time.Now()
	file, err := os.Open("./airline.csv.sample")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Insert into db
	batch := session.NewBatch(gocql.LoggedBatch)
	stmt := `INSERT INTO demo.flight (uniquecarrier, flightnum, tailnum, actualelapsedtime, crselapsedtime, distance) VALUES (?, ?, ?, ?, ?, ?);`
	// encoding := "iso8859-1"
	// Create a scanner to scan file
	count := 0
	batchCount := 0
	batchSize := 10

	// Create a reader with ISO 8859-1 encoding
	reader := bufio.NewReader(transform.NewReader(file, charmap.ISO8859_1.NewDecoder()))
	scanner := bufio.NewScanner(reader)
	// scanner := bufio.NewScanner(file)

	// For each line of the file
	for scanner.Scan() {

		line := strings.Split(scanner.Text(), ",")

		actualelapsedtime := line[0]
		// airtime := line[1]
		// arrdelay := line[2]
		// arrtime := line[3]
		// crsarrtime := line[4]
		// crsdeptime := line[5]
		crselapsedtime := line[6]
		// cancellationcode := line[7]
		// cancelled := line[8]
		// carrierdelay := line[9]
		// dayofweek := line[10]
		// dayofmonth := line[11]
		// depdelay := line[12]
		// deptime := line[13]
		// dest := line[14]
		distance := line[15]
		// diverted := line[16]
		flightnum := line[17]
		// lateaircraftdelay := line[18]
		// month := line[19]
		// nasdelay := line[20]
		// origin := line[21]
		// securitydelay := line[22]
		tailnum := line[23]
		// taxiin := line[24]
		// taxiout := line[25]
		uniquecarrier := line[26]
		// weatherdelay := line[27]
		// year := line[28]

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

	err = session.ExecuteBatch(batch)
	if err != nil {
		log.Panic(err)
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	log.Printf("Elapsed time %s", time.Since(start))

}
