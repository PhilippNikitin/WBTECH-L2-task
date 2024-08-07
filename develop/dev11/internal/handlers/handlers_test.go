package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"dev11/internal/db"
	"dev11/internal/models"
)

func TestCreateEventHandler(t *testing.T) {
	// Setup
	db.Events = []models.Event{}
	currentID = 0 // Reset the currentID to ensure tests start with a known state

	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-08-06")
	form.Add("title", "Test Event")

	req := httptest.NewRequest("POST", "/create_event", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateEventHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"result":"event created"}`
	got := strings.TrimSpace(rr.Body.String())
	if got != expected {
		t.Errorf("handler returned unexpected body:\n got: %s\nwant: %s", got, expected)
	}

	// Check if event was created
	if len(db.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(db.Events))
	}

	event := db.Events[0]
	expectedDate := time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC)
	if event.UserID != 1 || event.Date != expectedDate || event.Title != "Test Event" {
		t.Errorf("event fields do not match:\n got: %+v\nwant: %+v", event, models.Event{
			UserID: 1,
			Date:   expectedDate,
			Title:  "Test Event",
		})
	}
}

func TestUpdateEventHandler(t *testing.T) {
	// Setup initial event
	db.Events = []models.Event{
		{ID: 1, UserID: 1, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Old Title"},
	}
	form := url.Values{}
	form.Add("id", "1")
	form.Add("user_id", "1")
	form.Add("date", "2024-08-06")
	form.Add("title", "Updated Title")

	req := httptest.NewRequest("POST", "/update_event", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UpdateEventHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"result":"event updated"}`
	got := strings.TrimSpace(rr.Body.String())
	if got != expected {
		t.Errorf("handler returned unexpected body:\n got: %s\nwant: %s", got, expected)
	}

	// Check if event was updated
	if len(db.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(db.Events))
	}

	event := db.Events[0]
	if event.Title != "Updated Title" {
		t.Errorf("event title was not updated")
	}
}

func TestDeleteEventHandler(t *testing.T) {
	// Setup initial event
	db.Events = []models.Event{
		{ID: 1, UserID: 1, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Event to Delete"},
	}
	form := url.Values{}
	form.Add("id", "1")

	req := httptest.NewRequest("POST", "/delete_event", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(DeleteEventHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"result":"event deleted"}`
	got := strings.TrimSpace(rr.Body.String())
	if got != expected {
		t.Errorf("handler returned unexpected body:\n got: %s\nwant: %s", got, expected)
	}

	// Check if event was deleted
	if len(db.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(db.Events))
	}
}

func TestEventsForDayHandler(t *testing.T) {
	// Setup initial events
	db.Events = []models.Event{
		{ID: 1, UserID: 1, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Event 1"},
		{ID: 2, UserID: 1, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Event 2"},
		{ID: 3, UserID: 1, Date: time.Date(2024, time.August, 7, 0, 0, 0, 0, time.UTC), Title: "Event 3"},
		{ID: 4, UserID: 2, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Event 4"},
	}
	req := httptest.NewRequest("GET", "/events_for_day?user_id=1&date=2024-08-06", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(EventsForDayHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `[{"id":1,"user_id":1,"date":"2024-08-06T00:00:00Z","title":"Event 1"},{"id":2,"user_id":1,"date":"2024-08-06T00:00:00Z","title":"Event 2"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestEventsForWeekHandler(t *testing.T) {
	// Setup initial events
	db.Events = []models.Event{
		{ID: 1, UserID: 1, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Event 1"},
		{ID: 2, UserID: 1, Date: time.Date(2024, time.August, 7, 0, 0, 0, 0, time.UTC), Title: "Event 2"},
		{ID: 3, UserID: 1, Date: time.Date(2024, time.August, 8, 0, 0, 0, 0, time.UTC), Title: "Event 3"},
		{ID: 4, UserID: 1, Date: time.Date(2024, time.August, 9, 0, 0, 0, 0, time.UTC), Title: "Event 4"},
	}
	req := httptest.NewRequest("GET", "/events_for_week?user_id=1&date=2024-08-06", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(EventsForWeekHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `[{"id":1,"user_id":1,"date":"2024-08-06T00:00:00Z","title":"Event 1"},{"id":2,"user_id":1,"date":"2024-08-07T00:00:00Z","title":"Event 2"},{"id":3,"user_id":1,"date":"2024-08-08T00:00:00Z","title":"Event 3"},{"id":4,"user_id":1,"date":"2024-08-09T00:00:00Z","title":"Event 4"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestEventsForMonthHandler(t *testing.T) {
	// Setup initial events
	db.Events = []models.Event{
		{ID: 1, UserID: 1, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Event 1"},
		{ID: 2, UserID: 1, Date: time.Date(2024, time.August, 7, 0, 0, 0, 0, time.UTC), Title: "Event 2"},
		{ID: 3, UserID: 1, Date: time.Date(2024, time.August, 6, 0, 0, 0, 0, time.UTC), Title: "Event 3"},
	}
	req := httptest.NewRequest("GET", "/events_for_month?user_id=1&date=2024-08-01", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(EventsForMonthHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `[{"id":1,"user_id":1,"date":"2024-08-06T00:00:00Z","title":"Event 1"},{"id":2,"user_id":1,"date":"2024-08-07T00:00:00Z","title":"Event 2"},{"id":3,"user_id":1,"date":"2024-08-06T00:00:00Z","title":"Event 3"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
