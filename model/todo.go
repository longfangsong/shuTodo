package model

import (
	"shuTodo/infrastructure"
	"time"
)

type Todo struct {
	Id           int64         `json:"id"`
	Content      string        `json:"content"`
	Due          time.Time     `json:"due"`
	EstimateCost time.Duration `json:"estimate_cost"`
	Type         string        `json:"type"`
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
	err := row.Scan(&object.Content, &object.Due, &object.EstimateCost, &object.Type)
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
		err = rows.Scan(&item.Id, &item.Content, &item.Due, &item.EstimateCost, &item.Type)
		if err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func SaveTodo(object Todo) (Todo, error) {
	if object.Id == 0 {
		row := infrastructure.DB.QueryRow(`
		INSERT INTO Todo(content, due, estimatecost, type)
		VALUES ($1, $2, $3, $4)
		returning id;
		`, object.Content, object.Due, object.EstimateCost, object.Type)
		err := row.Scan(&object.Id)
		return object, err
	} else {
		_, err := infrastructure.DB.Exec(`
		UPDATE Todo
		SET content=$2,
		    due=$3,
		    estimatecost=$4,
		    type=$5
		WHERE id=$1;
		`, object.Id, object.Content, object.Due, object.EstimateCost, object.Type)
		return object, err
	}
}

func AssignTodoToStudent(studentId string, todoId int64) error {
	_, err := infrastructure.DB.Exec(`
	INSERT INTO studenttodo(student_id, todo_id)
	VALUES ($1, $2)
	ON CONFLICT DO UPDATE set student_id=$1;
	`, studentId, todoId)
	return err
}
