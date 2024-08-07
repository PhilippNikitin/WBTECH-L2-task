package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"dev11/internal/models"
)

// Helper function to write JSON responses
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// Helper function to parse form and validate parameters
func ParseFormAndValidate(r *http.Request) (models.Event, error) {
	var event models.Event
	if err := r.ParseForm(); err != nil {
		return event, err
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return event, fmt.Errorf("invalid user_id: %v", err)
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		return event, fmt.Errorf("invalid date: %v", err)
	}

	event.UserID = userID
	event.Date = date
	event.Title = r.FormValue("title")
	return event, nil
}

// Helper function to parse and validate GET parameters
func ParseGetParams(r *http.Request) (int, time.Time, error) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("invalid user_id: %v", err)
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("invalid date: %v", err)
	}

	return userID, date, nil
}

// Helper function to filter events by date range and user ID
func FilterEvents(events []models.Event, userID int, startDate, endDate time.Time) []models.Event {
	var filteredEvents []models.Event
	for _, e := range events {
		if e.UserID == userID && (e.Date.Equal(startDate) || (e.Date.After(startDate) && e.Date.Before(endDate))) {
			filteredEvents = append(filteredEvents, e)
		}
	}
	return filteredEvents
}
