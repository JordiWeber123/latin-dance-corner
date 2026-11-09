package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
	req := httptest.NewRequest(http.MethodGet, "/api/v1/move/1", nil)
	rec := httptest.NewRecorder()
	
	apiCfg := ApiConfig{DB: nil}
	
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

	testMove := Move{
		Name: "Inside turn",
		DanceStyle: "Bachata",
		Difficulty: 1,
		Description: "Lead: Prepare a hand by raising it up between counts 3 and 4, then on 5 pass it over and around the follow's head, forming a circle",
		Image: "img.png",}
		
	if body.Move != testMove{
		t.Errorf("message = %v, want %v", body.Move, testMove)
	}
}