package main_test

import (
	"."
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	db *main.DBServer
	id string
)

func testMain(m *testing.M) {
	connectDB()
	code := m.Run()
	os.Exit(code)
}

func connectDB() {
	db = &main.DBServer{Server: "localhost:27017", Database: "address_book"}
	db.Connect()
}

func TestNotFoundResponse(t *testing.T) {
	db.DeleteCollection()

	req, _ := http.NewRequest("GET", "/peoples", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(main.GetAllPeoples)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
	}
}

func TestCreatePeople(t *testing.T) {
	db.DeleteCollection()

	data := []byte(`{"fullname":"batman","phone":"742671","age":"20"}`)
	req, _ := http.NewRequest("POST", "/peoples", bytes.NewBuffer(data))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(main.CreatePeople)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}

	var p map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &p)

	if p["fullname"] != "batman" {
		t.Errorf("Expected fullname to be 'batman' got %v", p["fullname"])
	}
	if p["age"] != "20" {
		t.Errorf("Expected age to be '20' got %v", p["age"])
	}
	id = p["id"].(string)
}

func TestListPeoples(t *testing.T) {
	req, _ := http.NewRequest("GET", "/peoples", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(main.GetAllPeoples)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestGetPeople(t *testing.T) {
	req, _ := http.NewRequest("GET", "/peoples/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(main.GetPeople)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var p map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &p)

	if p["fullname"] != "batman" {
		t.Errorf("Expected fullname to be 'batman' got %v", p["fullname"])
	}
}

func TestUpdatePeople(t *testing.T) {
	payload := fmt.Sprintf(`{"id":"%s","fullname":"robin","phone":"742671","age":"20"}`, id)
	data := []byte(payload)
	req, _ := http.NewRequest("PUT", "/peoples", bytes.NewBuffer(data))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(main.UpdatePeople)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestDeletePeople(t *testing.T) {
	payload := fmt.Sprintf(`{"id":"%s","fullname":"robin","phone":"742671","age":"20"}`, id)
	data := []byte(payload)
	req, _ := http.NewRequest("DELETE", "/peoples", bytes.NewBuffer(data))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(main.DeletePeople)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestNotFoundResponse2(t *testing.T) {
	db.DeleteCollection()

	req, _ := http.NewRequest("GET", "/peoples", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(main.GetAllPeoples)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
	}
}
