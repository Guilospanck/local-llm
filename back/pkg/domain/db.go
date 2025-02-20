package domain

import (
	"base/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	hostname string
	dbname   string
	db       *sql.DB
}

func (mydb *Database) Connect() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", utils.DB_USER, utils.DB_PASSWORD, mydb.hostname, mydb.dbname)
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

func (mydb *Database) getPropertyIdFromViewIds(viewIds string) []int {
	queriedPropertyIds := []int{}

	query := fmt.Sprintf(Queries["property_views"]["id_viewid"], viewIds)
	stmt, err := mydb.db.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing query %s:\n %s\n", Queries["property_views"]["id_viewid"], err.Error())
		panic("")
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("\nError executing query %s:\n %s\n", Queries["property_views"]["id_viewid"], err.Error())
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

	query := fmt.Sprintf(Queries["view"]["id_view"], "'"+view+"%'")

	stmt, err := mydb.db.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing query %s:\n %s\n", Queries["view"]["id_view"], err.Error())
		panic(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("\nError executing query %s:\n %s\n", Queries["view"]["id_view"], err.Error())
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

		viewIds := utils.JoinSlice(queriedViewIds, ",")

		// Get property ids based on the views ids
		if viewIds != "" {
			propertyIds := mydb.getPropertyIdFromViewIds(viewIds)
			queriedPropertyIds = append(queriedPropertyIds, propertyIds...)
		}
	}

	propertyIds := utils.JoinSlice(queriedPropertyIds, ",")

	queryByColor := fmt.Sprintf(Queries["property"]["*_color"], color)
	queryByPrice := fmt.Sprintf(Queries["property"]["*_price"], priceMin, priceMax)
	queryBySize := fmt.Sprintf(Queries["property"]["*_size_sqm"], sizeSqmMin, sizeSqmMax)

	query := fmt.Sprintf("%s UNION %s UNION %s", queryByColor, queryByPrice, queryBySize)

	if propertyIds != "" {
		queryPropertyByIds := fmt.Sprintf(Queries["property"]["*_id"], propertyIds)
		query += fmt.Sprintf(" UNION %s", queryPropertyByIds)
	}

	if limit {
		query += fmt.Sprintf(" LIMIT %d;", utils.MAX_ITEMS_TO_QUERY)
	}

	fmt.Printf("\n\nQuerying DB:\n %s\n\n", query)

	stmt, err := mydb.db.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing query %s:\n %s\n", query, err.Error())
		panic(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("\nError executing query %s:\n %s\n", query, err.Error())
		panic(err)
	}
	defer rows.Close()

	properties := []Property{}

	for rows.Next() {
		property := Property{}
		err = rows.Scan(&property.Id, &property.Color, &property.Price, &property.SizeSqm)
		if err != nil {
			fmt.Printf("\nError scanning query %s:\n %s\n", query, err.Error())
			panic(err)
		}
		properties = append(properties, property)
	}

	return properties
}

func (mydb *Database) QueryAll() []Property {
	rows, err := mydb.db.Query(Queries["property"]["*"])
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	properties := []Property{}
	for rows.Next() {
		property := Property{}
		err = rows.Scan(&property.Id, &property.Color, &property.Price, &property.SizeSqm)
		if err != nil {
			fmt.Printf("\nError scanning query %s:\n %s\n", Queries["property"]["*"], err.Error())
			panic(err)
		}
		properties = append(properties, property)
	}

	return properties
}

func NewDb() *Database {
	hostname, exists := os.LookupEnv("DB_HOSTNAME")
	if !exists {
		hostname = utils.DB_HOSTNAME
	}

	return &Database{
		hostname: hostname,
		dbname:   utils.DB_DBNAME,
	}
}
