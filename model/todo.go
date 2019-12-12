package model

import (
	"database/sql"
	"errors"
	"shuTodo/infrastructure"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	Id           int64          `json:"id"`
	Content      string         `json:"content"`
	Due          *time.Time     `json:"due,omitempty"`
	EstimateCost *time.Duration `json:"estimate_cost,omitempty"`
	Type         string         `json:"type,omitempty"`
}

func GetTodo(id int64) (Todo, error) {
	row := infrastructure.DB.QueryRow(`
	SELECT content, due, estimatecost, type
	FROM Todo
	WHERE id=$1;
	`, id)
	object := Todo{
		Id: id,
	}
	var estimateCost *string
	err := row.Scan(&object.Content, &object.Due, &estimateCost, &object.Type)
	if estimateCost != nil {
		hourMinuteSecond := strings.Split(*estimateCost, ":")
		hour, _ := strconv.Atoi(hourMinuteSecond[0])
		minute, _ := strconv.Atoi(hourMinuteSecond[1])
		cost := time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute
		object.EstimateCost = &cost
	}
	return object, err
}

func GetTodoByStudentId(studentId string) ([]Todo, error) {
	rows, err := infrastructure.DB.Query(`
	SELECT id, content, due, estimatecost, type
	FROM todo,studenttodo
	where id=studenttodo.todo_id AND studenttodo.student_id=$1;
	`, studentId)
	if err != nil {
		return nil, err
	}
	var result []Todo
	for rows.Next() {
		var item Todo
		var estimateCost *string
		err = rows.Scan(&item.Id, &item.Content, &item.Due, &estimateCost, &item.Type)
		if estimateCost != nil {
			hourMinuteSecond := strings.Split(*estimateCost, ":")
			hour, _ := strconv.Atoi(hourMinuteSecond[0])
			minute, _ := strconv.Atoi(hourMinuteSecond[1])
			cost := time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute
			item.EstimateCost = &cost
		}
		if err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func SaveTodo(object Todo) (Todo, error) {
	if object.Id == 0 {
		var estimateCost sql.NullString
		if object.EstimateCost != nil {
			estimateCost.String = object.EstimateCost.String()
			estimateCost.Valid = true
		} else {
			estimateCost.Valid = false
		}
		var dueTime sql.NullTime
		if object.Due != nil {
			dueTime.Time = *object.Due
			dueTime.Valid = true
		} else {
			dueTime.Valid = false
		}
		row := infrastructure.DB.QueryRow(`
		INSERT INTO Todo(content, due, estimatecost, type)
		VALUES ($1, $2, $3, $4)
		returning id;
		`, object.Content, object.Due, estimateCost, object.Type)
		err := row.Scan(&object.Id)
		return object, err
	} else {
		var estimateCost sql.NullString
		if object.EstimateCost != nil {
			estimateCost.String = object.EstimateCost.String()
			estimateCost.Valid = true
		} else {
			estimateCost.Valid = false
		}
		_, err := infrastructure.DB.Exec(`
		UPDATE Todo
		SET content=$2,
		    due=$3,
		    estimatecost=$4,
		    type=$5
		WHERE id=$1;
		`, object.Id, object.Content, object.Due, estimateCost, object.Type)
		return object, err
	}
}

func AssignTodoToStudent(studentId string, todoId int64) error {
	_, err := infrastructure.DB.Exec(`
	INSERT INTO studenttodo(student_id, todo_id)
	VALUES ($1, $2)
	ON CONFLICT(todo_id) DO UPDATE set student_id=$1;
	`, studentId, todoId)
	return err
}

func DeleteTodoByStudent(studentId string, todoId int64) error {
	result, err := infrastructure.DB.Exec(`
	DELETE FROM Todo
	WHERE id = $2
  		AND id in (SELECT todo_id
             FROM studenttodo
             where student_id = $1);
	`, studentId, todoId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return err
}
