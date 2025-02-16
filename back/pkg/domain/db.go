package domain

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	_ "github.com/lib/pq"
)

type Database struct {
	hostname string
	dbname   string
	db       *sql.DB
}

const DB_USER = "postgres"
const DB_PASSWORD = "postgres"

func (mydb *Database) Connect() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, mydb.hostname, mydb.dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		panic(err)
	}

	fmt.Printf("%s\n", connStr)
	fmt.Println("Connected to Postgres!")

	mydb.db = db
}

func (mydb *Database) Close() {
	mydb.db.Close()
}

type Property struct {
	Id      int     `json:"-"`
	Color   string  `json:"color"`
	Price   float64 `json:"price"`
	SizeSqm float32 `json:"sizeSqm"`
}

func (mydb *Database) QueryByCharacteristics(color string, priceMin, priceMax float64, sizeSqmMin, sizeSqmMax float32, limit bool) []Property {
	if sizeSqmMax == 0 {
		sizeSqmMax = math.MaxFloat32
	}
	if priceMax == 0 {
		priceMax = math.MaxFloat64
	}

	queryLimit := ""
	if limit {
		queryLimit = "LIMIT 3"
	}

	stmt, err := mydb.db.Prepare(fmt.Sprintf(`
		SELECT * FROM property
			WHERE color = $1 UNION
		SELECT * FROM property
			WHERE price >= $2 AND price <= $3 UNION
		SELECT * FROM property
			WHERE size_sqm >= $4 AND size_sqm <= $5 %s;
		`, queryLimit))

	if err != nil {
		fmt.Printf("Error preparing query: %s", err.Error())
		panic(err)
	}

	rows, err := stmt.Query(color, priceMin, priceMax, sizeSqmMin, sizeSqmMax)
	if err != nil {
		fmt.Printf("Error executing query: %s", err.Error())
		panic(err)
	}
	defer rows.Close()

	properties := []Property{}

	for rows.Next() {
		property := Property{}
		err = rows.Scan(&property.Id, &property.Color, &property.Price, &property.SizeSqm)
		if err != nil {
			// handle this error
			panic(err)
		}
		properties = append(properties, property)
	}

	return properties
}

func (mydb *Database) QueryAll() []Property {
	rows, err := mydb.db.Query("SELECT * FROM property")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	properties := []Property{}
	for rows.Next() {
		property := Property{}
		err = rows.Scan(&property.Id, &property.Color, &property.Price, &property.SizeSqm)
		if err != nil {
			// handle this error
			panic(err)
		}
		properties = append(properties, property)
	}

	return properties
}

func NewDb() *Database {
	return &Database{
		hostname: "localhost",
		dbname:   "local-ai",
	}
}
