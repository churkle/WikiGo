// +build integration

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

func dbURLFromResource(r *dockertest.Resource) string {
	port := r.GetPort("5432/tcp")
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"user",
		"password",
		"localhost",
		port,
		"project")
}

func TestCreateDB(m *testing.M) {
	pool, err := dockertest.NewPool("")
	var db *sql.DB

	if err != nil {
		t.Errorf("Could not connect to docker: %s", err)
	}

	//Pulls image and
	opts := dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "latest",
		Env:          []string{"POSTGRES_PASSWORD=password", "POSTGRES_USER=user", "POSTGRES_DB=project"},
		ExposedPorts: []string{"5432"},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		t.Errorf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		db, err = sql.Open("postgres", dbURLFromResource(resource))
		if err != nil {
			log.Println("Database not ready yet (it is booting up, wait for a few tries)...")
			return err
		}

		// Tests if database is reachable
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	code := m.Run()

	// Delete the Docker container
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
