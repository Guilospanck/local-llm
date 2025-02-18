package domain

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

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

func joinSlice[T []E, E any](slice T, separator string) string {
	strArray := make([]string, len(slice))
	for i, element := range slice {
		strArray[i] = fmt.Sprint(element)
	}

	return strings.Join(strArray, separator)
}

func (mydb *Database) getPropertyIdFromViewIds(viewIds string) []int {
	queriedPropertyIds := []int{}

	query := fmt.Sprintf(`SELECT DISTINCT property_id from property_views WHERE view_id in (%s)`, viewIds)
	stmt, err := mydb.db.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing query SELECT property_id from property_views:\n %s\n", err.Error())
		panic("")
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("\nError executing query SELECT property_id FROM property_views:\n %s\n", err.Error())
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var propertyId int
		err = rows.Scan(&propertyId)
		if err != nil {
			// handle this error
			panic(err)
		}
		queriedPropertyIds = append(queriedPropertyIds, propertyId)
	}

	return queriedPropertyIds
}

func (mydb *Database) getViewIds(view string) []int {
	queriedViewIds := []int{}

	query := fmt.Sprintf(`SELECT v.id from "view" v WHERE v.view ILIKE %s`, "'"+view+"%'")

	stmt, err := mydb.db.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing query SELECT id from view:\n %s\n", err.Error())
		panic(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("\nError executing query SELECT id from VIEW:\n %s\n", err.Error())
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var viewId int
		err = rows.Scan(&viewId)
		if err != nil {
			// handle this error
			panic(err)
		}
		queriedViewIds = append(queriedViewIds, viewId)
	}

	return queriedViewIds
}

func (mydb *Database) QueryByCharacteristics(color string, priceMin, priceMax float64, sizeSqmMin, sizeSqmMax float32, views []string, limit bool) []Property {
	if sizeSqmMax == 0 {
		sizeSqmMax = math.MaxFloat32
	}
	if priceMax == 0 {
		priceMax = math.MaxFloat64
	}

	queriedPropertyIds := []int{}
	if len(views) > 0 {
		queriedViewIds := []int{}

		// Get view table ids based on the `views` slice
		for _, queryView := range views {
			viewIds := mydb.getViewIds(queryView)
			queriedViewIds = append(queriedViewIds, viewIds...)
		}

		viewIds := joinSlice(queriedViewIds, ",")

		// Get property ids based on the views ids
		if viewIds != "" {
			propertyIds := mydb.getPropertyIdFromViewIds(viewIds)
			queriedPropertyIds = append(queriedPropertyIds, propertyIds...)
		}
	}

	propertyIds := joinSlice(queriedPropertyIds, ",")

	queryByColor := fmt.Sprintf("SELECT * FROM property WHERE color = '%s'", color)
	queryByPrice := fmt.Sprintf("SELECT * FROM property WHERE price >= %f AND price <= %f", priceMin, priceMax)
	queryBySize := fmt.Sprintf("SELECT * FROM property WHERE size_sqm >= %f AND size_sqm <= %f", sizeSqmMin, sizeSqmMax)

	query := fmt.Sprintf("%s UNION %s UNION %s", queryByColor, queryByPrice, queryBySize)

	if propertyIds != "" {
		query += fmt.Sprintf(" UNION SELECT * FROM property WHERE id in (%s)", propertyIds)
	}

	if limit {
		query += " LIMIT 3;"
	}

	fmt.Printf("\n\nQuerying DB:\n %s\n\n", query)

	stmt, err := mydb.db.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing query: %s", err.Error())
		panic(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("\nError executing query SELECT * FROM property:\n %s\n", err.Error())
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
