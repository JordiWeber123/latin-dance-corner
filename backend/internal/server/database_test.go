package server

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

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

func TestDBNegativeId(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	testId := int64(-1)
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", TABLE_NAME)
	mock.ExpectQuery(query).
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

func TestDBCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	testMove := Move{
		Name: "Inside turn",
		DanceStyle: "Bachata",
		Difficulty: 1,
		Description: "Lead: Prepare a hand by raising it up between counts 3 and 4, then on 5 pass it over and around the follow's head, forming a circle",
		Image: "img.png",
	}
	query := fmt.Sprintf("INSERT INTO %s (name, dance_style, difficulty, description, image) VALUES ($1, $2, $3, $4, $5)", TABLE_NAME)
	mock.ExpectExec(query).
		WithArgs(testMove.Name,
			testMove.DanceStyle,
			testMove.Difficulty,
			testMove.Description,
			testMove.Image,
		).WillReturnResult(sqlmock.NewResult(1,1))
	res, err := CreateMove(db, testMove)
	if err != nil {
		t.Fatalf("could not create move: %v", err)
	}
	//Actually, probably want to return id too=
	if res != 1 {
		t.Fatalf("did not create 1 row, created %d instead", res)
	}
}

func TestDBDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	testId := 1
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", TABLE_NAME)
	mock.ExpectExec(query).
		WithArgs(testId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	res, err := DeleteMove(db, int64(testId))
	
	if err != nil {
		t.Fatal(err)
	}

	if !res {
		t.Fatalf("could not delete move with id=%v", testId)
	} 
}

func TestDBDeleteNotFound(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	testId := 1
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", TABLE_NAME)
	mock.ExpectExec(query).
		WithArgs(testId).
		WillReturnResult(sqlmock.NewResult(1, 0))
	res, err := DeleteMove(db, int64(testId))
	
	if err != nil {
		t.Fatal(err)
	}

	if res {
		t.Fatal("returned true when expected false")
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


func TestDBMulti(t *testing.T) {
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
	mockRows := sqlmock.NewRows([]string{
			"id", "name", "dance_style", "difficulty", 
			"description", "image", "created_at",
	})

	for i := int64(0); i < 10; i++ {
		mockRows.AddRow(
			testMove.Id + i,
			testMove.Name + strconv.FormatInt(i, 10),
			testMove.DanceStyle,
			testMove.Difficulty,
			testMove.Description,
			testMove.Image,
			testMove.CreatedAt.Add(time.Duration(i) * time.Second),
		)
	}
	testPageSize := uint(10)
	testSkip := uint(0)

	mock.ExpectQuery("SELECT * FROM moves LIMIT $1 OFFSET $2").
		WithArgs(testPageSize, testSkip).
		WillReturnRows(mockRows)
	
	moves, err := GetMultiPaginated(db, testPageSize, testSkip)
	if err != nil {
		t.Fatalf("there was an error calling GetMultiPaginated: %v", err)
	}

	if len(moves) != int(testPageSize) {
		t.Fatalf("length of moves is: %d, want %d", len(moves), int(testPageSize))

	}
}

func TestDBMultiEmpty(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	mockRows := sqlmock.NewRows([]string{
			"id", "name", "dance_style", "difficulty", 
			"description", "image", "created_at",
	})
	testPageSize := uint(10)
	testSkip := uint(0)

	mock.ExpectQuery("SELECT * FROM moves LIMIT $1 OFFSET $2").
		WithArgs(testPageSize, testSkip).
		WillReturnRows(mockRows)
	
	moves, err := GetMultiPaginated(db, testPageSize, testSkip)
	if err != nil {
		t.Fatalf("there was an error calling GetMultiPaginated: %v", err)
	}

	if len(moves) != 0 {
		t.Fatalf("length of moves is: %d, want %d", len(moves), 0)

	}
}

func TestDBLimitZero(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()
	mockRows := sqlmock.NewRows([]string{
			"id", "name", "dance_style", "difficulty", 
			"description", "image", "created_at",
	})
	testPageSize := uint(0)
	testSkip := uint(0)

	mock.ExpectQuery("SELECT * FROM moves LIMIT $1 OFFSET $2").
		WithArgs(testPageSize, testSkip).
		WillReturnRows(mockRows)
	
	moves, err := GetMultiPaginated(db, testPageSize, testSkip)
	if err != nil {
		t.Fatalf("there was an error calling GetMultiPaginated: %v", err)
	}

	if len(moves) != int(testPageSize) {
		t.Fatalf("length of moves is: %d, want %d", len(moves), testPageSize)
	}
}

func TestDBLessPageSize(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("no database: %v", err)
	}
	defer db.Close()

	testMove := Move{
		Id: 1,
		Name: "Inside turn",
		DanceStyle: "Bachata",
		Difficulty: 1,
		Description: "Lead: Prepare a hand by raising it up between counts 3 and 4, then on 5 pass it over and around the follow's head, forming a circle",
		Image: "img.png",
	}
	mockRows := sqlmock.NewRows([]string{
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
	)
	testPageSize := uint(10)
	testSkip := uint(0)

	mock.ExpectQuery("SELECT * FROM moves LIMIT $1 OFFSET $2").
		WithArgs(testPageSize, testSkip).
		WillReturnRows(mockRows)
	
	moves, err := GetMultiPaginated(db, testPageSize, testSkip)
	if err != nil {
		t.Fatalf("there was an error calling GetMultiPaginated: %v", err)
	}

	if len(moves) > int(testPageSize) {
		t.Fatalf("length of moves is: %d, want <=%d", len(moves), testPageSize)
	}
}