package server

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var TEST_DB_CONN string = "postgres://testuser:testpass@127.0.0.1:5433/app_test?sslmode=disable"

func TestDBFailGetId(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	testId := int64(10)
	mock.ExpectQuery("SELECT * FROM moves WHERE id = $1").
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "dance_style", "difficulty", 
			"description", "image", "created_at",
		}))
	_, err = GetMoveById(db, testId)
	
	if !errors.Is(err, ErrMoveNotFound{}){
		t.Fatalf("didn't get ErrMoveNotFound")
	}
}

func TestDBGetId(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	testId := int64(1)
	testMove := Move{
		Id: testId,
		Name: "Inside turn",
		DanceStyle: "Bachata",
		Difficulty: 1,
		Description: "Lead: Prepare a hand by raising it up between counts 3 and 4, then on 5 pass it over and around the follow's head, forming a circle",
		Image: "img.png",
	}
	mock.ExpectQuery("SELECT * FROM moves WHERE id = $1").
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "dance_style", "difficulty", 
			"description", "image", "created_at",
		}).AddRow(
			testMove.Id,
			testMove.Name,
			testMove.DanceStyle,
			testMove.Difficulty,
			testMove.Description,
			testMove.Image,
			testMove.CreatedAt,
		))

	move, err := GetMoveById(db, testId)
	
	if err != nil {
		t.Fatalf("go unexpected error during test: %v", err)
	}

	if move != testMove {
		t.Fatalf("got move= %v want move=%v", move, testMove)
	}
}

//Integration tests:

func TestDBStoreRetrieve(t *testing.T){
	sql.Open("pgx", "")
	/*No integration test for now
	
	db, err := sql.Open("pgx", TEST_DB_CONN)
	if err != nil {
		t.Fatalf("error while opening DB: %s", err.Error())
	}
	defer db.Close()

	toAdd := Move{
		Name: "Inside turn",
		DanceStyle: "Bachata",
		Difficulty: 1,
		Description: "Lead: Prepare a hand by raising it up between counts 3 and 4, then on 5 pass it over and around the follow's head, forming a circle",
		Image: "img.png",
	}
	
	affected, err := CreateMove(db, toAdd)
	if err != nil {
		t.Fatalf("error storing move: %s", err.Error())
	}

	if affected != 1 {
		t.Fatalf("wrong number of affected rows %d != 1", affected)
	}
	
	moves, err := GetMultiPaginated(db, 1, 0)
	if err != nil {
		t.Fatalf("error getting moves list: %s", err.Error())
	}
	found := false
	for _, move := range moves {
		if (move.Name == toAdd.Name && move.Description == toAdd.Description && 
		move.Image == toAdd.Image){
			found = true
			break
		}	
	}
	
	if !found {
		t.Fatalf("move was not added")
	}*/
}