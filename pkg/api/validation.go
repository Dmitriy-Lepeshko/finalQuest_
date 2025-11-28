package api

import (
	"errors"
	"time"

	"firstIteration/pkg/db"
)

const dateFormat = "20060102"

func isBeforeDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	if y1 < y2 {
		return true
	}
	if y1 > y2 {
		return false
	}
	if m1 < m2 {
		return true
	}
	if m1 > m2 {
		return false
	}
	return d1 < d2
}

func validateTask(task *db.Task) error {
	if task.Title == "" {
		return errors.New("title is required")
	}

	if task.Date != "" && task.Date != "today" {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			return errors.New("invalid date format, expected YYYYMMDD")
		}
	}

	if task.Repeat != "" {
		dummyDate := "20240101"
		if _, err := NextDate(time.Now(), dummyDate, task.Repeat); err != nil {
			return err
		}
	}

	return nil
}

func checkDate(task *db.Task) error {
	now := time.Now()
	nowStr := now.Format(dateFormat)

	if task.Date == "" || task.Date == "today" {
		task.Date = nowStr
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format, expected YYYYMMDD")
	}

	if isBeforeDate(t, now) {
		if task.Repeat == "" {
			task.Date = nowStr
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}
