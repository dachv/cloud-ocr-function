package cloudsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	db         *sql.DB
	dbHost     = os.Getenv("POSTGRES_HOST")
	dbPort     = os.Getenv("POSTGRES_PORT")
	dbUser     = os.Getenv("POSTGRES_USER")
	dbPassword = os.Getenv("POSTGRES_PASSWORD")
	dbName     = os.Getenv("POSTGRES_DATABASE")
	dsName     = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
)

func init() {
	var err error
	db, err = sql.Open("postgres", dsName)
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
	// Only allow 1 connection to the database to avoid overloading it.
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
}

func InsertDocumentMetadata(metadata map[string]string) int {
	rows, err := db.Query("SELECT current_timestamp AS ts")
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
	defer rows.Close()
	ts := ""
	rows.Next()
	if err := rows.Scan(&ts); err != nil {
		log.Printf("rows.Scan: %v", err)
	}
	fmt.Printf("TS: %v", ts)
	return 0
}
