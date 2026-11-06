package main

import (
	"fmt"
	"database/sql"
	"github.com/JordiWeber123/latin-dance-corner"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	//placeholder
	var TEST_DB_CONN string = "postgres://testuser:testpass@127.0.0.1:5433/app_test?sslmode=disable"
	db, err := sql.Open("pgx", TEST_DB_CONN)
	if err != nil {
		fmt.Printf("error while opening DB: %s", err.Error())
		return
	}

	move, err := backend.GetMoveById(db, 1)
	fmt.Println(move, err)
}