package main

import (
	"context"
	"log"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	start := time.Now()
	defer func(start time.Time) {
		log.Printf("Elapsed time %s", time.Since(start))
	}(start)

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
      user_id int,
      fname text,
      lname text,
      PRIMARY KEY((user_id))
    )`,
	).WithContext(ctx).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Insert into db
	err = session.Query(
		`INSERT INTO demo.flight(
      user_id,
      fname,
      lname
    ) values (?, ?, ?)`,
		1337,
		"firstname",
		"lastname",
	).WithContext(ctx).Exec()
	if err != nil {
		log.Fatal(err)
	}

}
