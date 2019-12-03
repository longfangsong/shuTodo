package handler

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"io/ioutil"
	"net/http/httptest"
	"shuTodo/infrastructure"
	"shuTodo/tools"
	"testing"
	"time"
)

func TestCreateTodoHandler(t *testing.T) {
	var dbmock sqlmock.Sqlmock
	var err error
	infrastructure.DB, dbmock, err = sqlmock.New()
	tools.CheckErr(err, "cannot create Mock")
	duration, _ := time.ParseDuration("2h")
	location, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Date(2019, 12, 1, 10, 43, 47, 0, location)
	dbmock.ExpectQuery(`INSERT INTO Todo`).
		WithArgs("test", date, duration, "Homework").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
	dbmock.ExpectExec(`INSERT INTO studenttodo`).
		WithArgs("17120238", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	requestContent := struct {
		Id           int64  `json:"id"`
		Content      string `json:"content"`
		Due          string `json:"due"`
		EstimateCost string `json:"estimate_cost"`
		Type         string `json:"type"`
	}{
		Id:           0,
		Content:      "test",
		Due:          "2019-12-01T02:43:47.000Z",
		EstimateCost: "2h",
		Type:         "Homework",
	}
	requestBody, _ := json.Marshal(requestContent)
	r := httptest.NewRequest("POST", "http://localhost:8000/todo", bytes.NewReader(requestBody))
	// this token is my(studentId: 17120238) token encoded by JWT_SECRET "test"
	// for test only, never used in production
	r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdHVkZW50SWQiOiIxNzEyMDIzOCJ9.1shlMZ014Rnzw7Z5iNxiL73dC2xQ0iiKIFTILsOME-I")
	w := httptest.NewRecorder()
	CreateTodoHandler(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		t.Error("Creating failed with error:\n" + string(body))
	}
}

func TestGetTodoHandler(t *testing.T) {
	var dbmock sqlmock.Sqlmock
	var err error
	infrastructure.DB, dbmock, err = sqlmock.New()
	tools.CheckErr(err, "cannot create mock DB")
	duration, _ := time.ParseDuration("2h")
	location, _ := time.LoadLocation("Local")
	date := time.Date(2019, 12, 1, 0, 0, 0, 0, location)
	rows := sqlmock.NewRows([]string{"id", "content", "due", "estimatecost", "type"})
	rows.AddRow(1, "test", date, duration, "Homework")
	dbmock.ExpectQuery(`SELECT (.+) FROM (.+) where (.+);`).
		WithArgs("17120238").WillReturnRows(rows)
	r := httptest.NewRequest("POST", "http://localhost:8000/todo", nil)
	// this token is my(studentId: 17120238) token encoded by JWT_SECRET "test"
	// for test only, never used in production
	r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdHVkZW50SWQiOiIxNzEyMDIzOCJ9.1shlMZ014Rnzw7Z5iNxiL73dC2xQ0iiKIFTILsOME-I")
	w := httptest.NewRecorder()
	GetTodoHandler(w, r)
	body, _ := ioutil.ReadAll(w.Body)
	var todos []struct {
		Id           int64  `json:"id"`
		Content      string `json:"content"`
		Due          string `json:"due"`
		EstimateCost int64  `json:"estimate_cost"`
		Type         string `json:"type"`
	}
	err = json.Unmarshal(body, &todos)
	if err != nil {
		t.Error(err)
	}
}
