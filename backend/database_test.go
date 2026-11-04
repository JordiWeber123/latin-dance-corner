package backend

import (
	"database/sql"
	"testing"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var TEST_DB_CONN string = "postgres://testuser:testpass@127.0.0.1:5433/app_test?sslmode=disable"


func TestStoreRetrieve(t *testing.T){
	db, err := sql.Open("pgx", TEST_DB_CONN)
	if err != nil {
		t.Fatalf("error while opening DB: %s", err.Error())
	}
	toAdd := Move{Name: "Inside turn", 
	Description: "Lead: Prepare a hand by raising it up between counts 3 and 4, then on 5 pass it over and around the follow's head, forming a circle",
	Image: "img.png"}
	
	affected, err := StoreMove(db, toAdd)
	if err != nil {
		t.Fatalf("error storing move: %s", err.Error())
	}

	if affected != 1 {
		t.Fatalf("wrong number of affected rows %d != 1", affected)
	}
	
	moves, err := GetNMovesSorted(db, 1, "id")
	if err != nil {
		t.Fatalf("error getting moves list: %s", err.Error())
	}
	
	if (moves[0].Name != toAdd.Name || moves[0].Description != toAdd.Description || 
		moves[0].Image != toAdd.Image){
			t.Fatalf("moves are not equal")
	}	
}