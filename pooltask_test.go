package pooltask

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleListTasks(t *testing.T) {
	//Create HTTP GET request
	req, err := http.NewRequest("GET", "http://localhost:4000/task", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("This will mark the test as failed but it will continue the execution")
	}
	//Record the HTTP Response (httptest)
	rec := httptest.NewRecorder()
	//Dispatch the HTTP Request
	HandleListTasks(rec, req)
	//Add Assertions on the HTTP Status code and the response
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", req.Response.StatusCode)
	}
}

func TestHandleCreateTask(t *testing.T) {
	//Create HTTP POST request
	var jsonB = make(map[string]int)
	jsonB["processMeForThisMuchSeconds"] = 15
	b, _ := json.Marshal(jsonB)
	req, err := http.NewRequest("POST", "http://localhost:4000/task", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Could not execute the request")
	}
	//Record the HTTP Response (httptest)
	rec := httptest.NewRecorder()
	//Dispatch the HTTP Request
	HandleCreateTask(rec, req)
	//Add Assertions on the HTTP Status code and the response
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v want %v", res.StatusCode, http.StatusOK)
	}

	var cr CreateResponse

	decoder := json.NewDecoder(io.Reader(res.Body))
	err = decoder.Decode(&cr)
	if err != nil {
		t.Errorf("Could not decode, due: %v", err)
	}
}
