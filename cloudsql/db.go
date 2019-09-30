package cloudsql

import (
	"database/sql"
	"encoding/json"
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

func ExistsWith(companyName string, formType string) (bool, int) {
	jsonClauseBytes, err := json.Marshal(map[string]string{"FormType": formType})
	if err != nil {
		log.Fatalf("Error marshalling json clause map: %v", err)
	}
	row := db.QueryRow(`select distinct on (version) version from document where company_name = $1
  						and metadata @> $2 order by version desc`, companyName, string(jsonClauseBytes))
	version := 0
	switch err = row.Scan(&version); err {
	case sql.ErrNoRows:
		return false, version
	case nil:
		return version > 0, version
	default:
		log.Fatalf("Error checking document existence: %v", err)
		return false, version
	}
}

func InsertDocument(objectName string, companyName string, version int, metadata map[string]string) int {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		log.Fatalf("Error marshalling document metadata to json: %v", err)
	}
	row := db.QueryRow(`INSERT INTO document(created_at, object_name, company_name, version, metadata) VALUES 
		(current_timestamp, $1, $2, $3, $4) RETURNING id`, objectName, companyName, version, string(metadataBytes))
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Fatalf("Error inserting document metadata to db: %v", err)
	}
	return id
}
