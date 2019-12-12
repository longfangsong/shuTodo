package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shuTodo/model"
	"shuTodo/service/token"
	"strconv"
	"time"
)

type todoInput struct {
	Id           int64  `json:"id"`
	Content      string `json:"content"`
	Due          string `json:"due"`
	EstimateCost string `json:"estimate_cost"`
	Type         string `json:"type"`
}

func parseInput(input todoInput) (model.Todo, error) {
	var duePointer *time.Time
	due, err := time.Parse("2006-01-02T15:04:05.999999999Z", input.Due)
	if err != nil {
		duePointer = nil
	} else {
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err)
		}
		due = due.In(location)
		duePointer = &due
	}
	var estimateCostPointer *time.Duration
	estimateCost, err := time.ParseDuration(input.EstimateCost)
	if err != nil {
		estimateCostPointer = nil
	} else {
		estimateCostPointer = &estimateCost
	}
	return model.Todo{
		Id:           input.Id,
		Content:      input.Content,
		Due:          duePointer,
		EstimateCost: estimateCostPointer,
		Type:         input.Type,
	}, nil
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	if len(tokenInHeader) <= 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	studentId := token.StudentIdForToken(tokenInHeader[len("Bearer "):])
	if studentId == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	var toCreateInput todoInput
	err = json.Unmarshal(body, &toCreateInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	toCreate, err := parseInput(toCreateInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	// todo: start a transaction?
	toCreate, err = model.SaveTodo(toCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.AssignTodoToStudent(studentId, toCreate.Id)
	// transaction end here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseBody, err := json.Marshal(toCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	if len(tokenInHeader) <= 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	studentId := token.StudentIdForToken(tokenInHeader[7:])
	if studentId == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	todos, err := model.GetTodoByStudentId(studentId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	body, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if len(todos) == 0 {
		_, err = w.Write([]byte("[]"))
	} else {
		_, err = w.Write(body)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	tokenInHeader := r.Header.Get("Authorization")
	if len(tokenInHeader) <= 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	studentId := token.StudentIdForToken(tokenInHeader[7:])
	if studentId == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	err = model.DeleteTodoByStudent(studentId, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetTodoHandler(w, r)
	case "POST":
		CreateTodoHandler(w, r)
	case "PUT":
		CreateTodoHandler(w, r)
	case "DELETE":
		DeleteTodoHandler(w, r)
	}
}
