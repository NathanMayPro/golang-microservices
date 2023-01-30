package main

import (
"fmt"
"database/sql"
"reflect"
"log"

_ "github.com/lib/pq"
)

var db *sql.DB

type DatbaseTableUser struct {
	Id	   int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Usertype string `json:"usertype"`
}

// function that execute a query and parse the result with any
// struct that implements the interface
func scan(rows *sql.Rows, result interface{}) error {
	// This line retrieves the reflect.Value of the result argument. 
	// The reflect.Value type represents a value in Go.
	ptr := reflect.ValueOf(result)

	// This line checks if the kind of the reflect.Value is a pointer using the Kind method. 
	// If the kind is not a pointer, the function returns an error.
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %T", result)
	}

	// This line dereferences the pointer to get the underlying value.
	val := ptr.Elem()
	// This line checks if the kind of the underlying value is a struct using the Kind method. 
	// If the kind is not a struct, the function returns an error.
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %T", result)
	}

	// This line retrieves the column names from the query result.
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(columns) != val.NumField() {
		return fmt.Errorf("struct has %d fields, query returned %d columns", val.NumField(), len(columns))
	}
	// This line creates a slice of interface{} values with a length equal to the number of columns.
	values := make([]interface{}, len(columns))
	// This line loops over the slice of values and sets each value to the address 
	// of the corresponding field in the struct using the Field method and the Addr method.
	for i := range values {
		values[i] = val.Field(i).Addr().Interface()
	}
	// This line calls the .Scan method on the query result with the slice of values, 
	// which assigns the values from the query result to the fields in the struct.
	return rows.Scan(values...)
}



// This function will make a connection to the database only once.
func main() {
	var err error

	connStr := "postgres://postgres:docker@localhost:5555/world?sslmode=disable"
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	// this will be printed in the terminal, confirming the connection to the database
	fmt.Println("The database is connected")

	// run a query to the database
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user DatbaseTableUser
		err := scan(rows, &user)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(user)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}