package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

var errInvalidInput = errors.New("invalid input")

// NextDate returns the next execution date based on the repeat rule.
// If inputs are invalid or rule is not supported, it returns an error.
func NextDate(nowStr, dateStr, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errInvalidInput
	}

	now, err := time.Parse(dateLayout, nowStr)
	if err != nil {
		return "", errInvalidInput
	}
	start, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return "", errInvalidInput
	}

	fields := strings.Fields(repeat)
	if len(fields) == 0 {
		return "", errInvalidInput
	}

	switch fields[0] {
	case "y":
		if len(fields) != 1 {
			return "", errInvalidInput
		}
		next := start.AddDate(1, 0, 0)
		for !next.After(now) {
			next = next.AddDate(1, 0, 0)
		}
		return next.Format(dateLayout), nil
	case "d":
		if len(fields) != 2 {
			return "", errInvalidInput
		}
		days, err := strconv.Atoi(fields[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errInvalidInput
		}
		next := start.AddDate(0, 0, days)
		for !next.After(now) {
			next = next.AddDate(0, 0, days)
		}
		return next.Format(dateLayout), nil
	default:
		return "", errInvalidInput
	}
}

// NextDateHandler handles GET /api/nextdate requests.
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	next, err := NextDate(nowStr, dateStr, repeat)
	if err != nil {
		return
	}
	_, _ = w.Write([]byte(next))
}
