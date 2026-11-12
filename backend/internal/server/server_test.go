package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	rec := httptest.NewRecorder()
	
	apiCfg := ApiConfig{DB: nil}
	
	New(&apiCfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	if body.Message != "Hello, world!" {
		t.Errorf("message = %q, want %q", body.Message, "Hello, world!")
	}
}

func TestGetId(t *testing.T){
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("unable to open database: %v", err.Error())
	}
	defer db.Close()
	
	testId := int64(1)
	testMove := Move{
		Id: int64(testId),
		Name: "Inside turn",
		DanceStyle: "Bachata",
		Difficulty: 1,
		Description: "Lead: Prepare a hand by raising it up between counts 3 and 4, then on 5 pass it over and around the follow's head, forming a circle",
		Image: "img.png",
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", TABLE_NAME)
	mock.ExpectQuery(query).
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
	
	testUrl := fmt.Sprintf("/api/v1/move/%d", testId)
	req := httptest.NewRequest(http.MethodGet, testUrl, nil)
	rec := httptest.NewRecorder()
	apiCfg := ApiConfig{DB: db}
	New(&apiCfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body struct {
		Move Move `json:"move"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
		
	if !(body.Move.Id == testId &&
		body.Move.Name == testMove.Name && 
		body.Move.DanceStyle == testMove.DanceStyle &&
		body.Move.Difficulty == testMove.Difficulty &&
		body.Move.Description == testMove.Description && 
		body.Move.Image == testMove.Image ){
		t.Errorf("message = %v, want %v", body.Move, testMove)
	}


}

func TestGetMissingId(t *testing.T){
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("unable to open database: %v", err.Error())
	}
	defer db.Close()
	testId := int64(42)
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", TABLE_NAME)
	mock.ExpectQuery(query).
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "dance_style", "difficulty", 
			"description", "image", "created_at",
		}))	
	
	testUrl := fmt.Sprintf("/api/v1/move/%d", testId)
	req := httptest.NewRequest(http.MethodGet, testUrl, nil)
	rec := httptest.NewRecorder()
	apiCfg := ApiConfig{DB: db}
	New(&apiCfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	var body struct {
		Error struct {
			Code 	string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}

	if body.Error.Code != "not_found" {
		t.Fatalf("code = %s want %s", body.Error.Code, "not_found")
	}
}

func TestNegativeId(t *testing.T) {
	testUrl := fmt.Sprintf("/api/v1/move/%d", -1)
	req := httptest.NewRequest(http.MethodGet, testUrl, nil)
	rec := httptest.NewRecorder()
	apiCfg := ApiConfig{DB: nil}
	New(&apiCfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	var body struct {
		Error struct {
			Code 	string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}

	if body.Error.Code != "invalid_id" {
		t.Fatalf("code = %s want %s", body.Error.Code, "invalid_id")
	}
}

func TestStrId(t *testing.T) {
	testUrl := fmt.Sprintf("/api/v1/move/%s", "blah")
	req := httptest.NewRequest(http.MethodGet, testUrl, nil)
	rec := httptest.NewRecorder()
	apiCfg := ApiConfig{DB: nil}
	New(&apiCfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	var body struct {
		Error struct {
			Code 	string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}

	if body.Error.Code != "invalid_id" {
		t.Fatalf("code = %s want %s", body.Error.Code, "invalid_id")
	}
}