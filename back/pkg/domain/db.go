package domain

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	hostname string
	dbname   string
}

const DB_USER = "postgres"
const DB_PASSWORD = "postgres"

func (mydb *Database) Connect() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, mydb.hostname, mydb.dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		db.Close()
		panic(err)
	}

	fmt.Printf("%s", connStr)
	fmt.Println("Connected to Postgres!")

	// Insert example
	// color := "blue"
	// price := 1_000_000.00
	// sizeSqm := 450.28
	// _, err = db.Query("INSERT INTO property(color, price, size_sqm) VALUES($1, $2, $3)", color, price, sizeSqm)

	type property struct {
		id      int
		color   string
		price   float64
		sizeSqm float32
	}

	rows, err := db.Query("SELECT * FROM property")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		test := property{}
		err = rows.Scan(&test.id, &test.color, &test.price, &test.sizeSqm)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Printf("%v\n", test)
	}
}

func NewDb() *Database {
	return &Database{
		hostname: "localhost",
		dbname:   "local-ai",
	}

}
