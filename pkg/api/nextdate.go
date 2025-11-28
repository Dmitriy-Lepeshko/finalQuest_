package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	if y1 > y2 {
		return true
	}
	if y1 < y2 {
		return false
	}
	if m1 > m2 {
		return true
	}
	if m1 < m2 {
		return false
	}
	return d1 > d2
}

func daysInMonth(t time.Time) int {
	year, month, _ := t.Date()
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	_, _, day := lastDay.Date()
	return day
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart: %w", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("d rule requires exactly one number")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("d rule: invalid number")
		}
		if days <= 0 || days > 400 {
			return "", errors.New("d rule: number must be between 1 and 400")
		}

		date := start
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("y rule takes no arguments")
		}

		date := start
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("w rule requires exactly one argument (comma-separated days)")
		}

		dayTokens := strings.Split(parts[1], ",")
		allowedDays := make(map[int]bool)
		for _, token := range dayTokens {
			d, err := strconv.Atoi(strings.TrimSpace(token))
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("w rule: days must be integers from 1 (Mon) to 7 (Sun)")
			}
			allowedDays[d] = true
		}

		date := start
		for {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			normalized := weekday
			if normalized == 0 {
				normalized = 7
			}

			if allowedDays[normalized] && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}

			if date.Year() > now.Year()+10 {
				return "", errors.New("w rule: no matching date found within 10 years")
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("m rule: expected 'm <days>' or 'm <days> <months>'")
		}

		dayTokens := strings.Split(parts[1], ",")
		allowedDays := make(map[int]bool)
		for _, token := range dayTokens {
			d, err := strconv.Atoi(strings.TrimSpace(token))
			if err != nil {
				return "", errors.New("m rule: invalid day")
			}
			if d < -2 || (d > 0 && d > 31) {
				return "", errors.New("m rule: day must be from 1 to 31, or -1, -2")
			}
			allowedDays[d] = true
		}

		allowedMonths := make(map[int]bool)
		if len(parts) == 3 {
			monthTokens := strings.Split(parts[2], ",")
			for _, token := range monthTokens {
				m, err := strconv.Atoi(strings.TrimSpace(token))
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("m rule: months must be from 1 to 12")
				}
				allowedMonths[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				allowedMonths[i] = true
			}
		}

		date := start
		for {
			date = date.AddDate(0, 0, 1)

			_, month, day := date.Date()
			monthNum := int(month)
			dayNum := day

			matchedDay := allowedDays[dayNum] ||
				(allowedDays[-1] && dayNum == daysInMonth(date)) ||
				(allowedDays[-2] && dayNum == daysInMonth(date)-1)

			if matchedDay && allowedMonths[monthNum] && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}

			if date.Year() > now.Year()+10 {
				return "", errors.New("m rule: no matching date found within 10 years")
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule: %s", parts[0])
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dateStr == "" || repeat == "" {
		http.Error(w, "missing 'date' or 'repeat' parameter", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		parsed, err := time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid 'now' format, expected YYYYMMDD", http.StatusBadRequest)
			return
		}
		now = parsed
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
