package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (string, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	task.ID = strconv.FormatInt(id, 10)
	return task.ID, nil
}

func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		ORDER BY date ASC 
		LIMIT ?
	`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func SearchByKeyword(keyword string, limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE UPPER(title) LIKE UPPER(?) OR UPPER(comment) LIKE UPPER(?)
		ORDER BY date ASC 
		LIMIT ?
	`
	rows, err := DB.Query(query, "%"+keyword+"%", "%"+keyword+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

func SearchByDate(dateStr string, limit int) ([]*Task, error) {
	if _, err := time.Parse("20060102", dateStr); err != nil {
		return []*Task{}, nil
	}

	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date = ?
		ORDER BY date ASC 
		LIMIT ?
	`
	rows, err := DB.Query(query, dateStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE id = ?
	`
	err := DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ? 
		WHERE id = ?
	`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task: %s", task.ID)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}
	return nil
}

func UpdateDate(nextDate, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, nextDate, id)
	if err != nil {
		return fmt.Errorf("failed to update date: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}
	return nil
}
