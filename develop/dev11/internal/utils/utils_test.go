package utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"dev11/internal/models"
)

// TestWriteJSON tests the WriteJSON function
func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"result": "success"}
	WriteJSON(rr, http.StatusOK, data)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check the body
	expectedBody := `{"result":"success"}`
	if got := strings.TrimSpace(rr.Body.String()); got != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expectedBody)
	}
}

// TestParseFormAndValidate tests the ParseFormAndValidate function
func TestParseFormAndValidate(t *testing.T) {
	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-08-06")
	form.Add("title", "Event Title")

	req := httptest.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	event, err := ParseFormAndValidate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the parsed event
	if event.UserID != 1 {
		t.Errorf("expected user_id 1, got %d", event.UserID)
	}
	expectedDate, _ := time.Parse("2006-01-02", "2024-08-06")
	if !event.Date.Equal(expectedDate) {
		t.Errorf("expected date 2024-08-06, got %s", event.Date)
	}
	if event.Title != "Event Title" {
		t.Errorf("expected title 'Event Title', got %s", event.Title)
	}
}

// TestParseGetParams tests the ParseGetParams function
func TestParseGetParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?user_id=1&date=2024-08-06", nil)

	userID, date, err := ParseGetParams(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the parsed user ID
	if userID != 1 {
		t.Errorf("expected user_id 1, got %d", userID)
	}

	// Check the parsed date
	expectedDate, _ := time.Parse("2006-01-02", "2024-08-06")
	if !date.Equal(expectedDate) {
		t.Errorf("expected date 2024-08-06, got %s", date)
	}
}

// TestFilterEvents tests the FilterEvents function
func TestFilterEvents(t *testing.T) {
	// Set up test data
	events := []models.Event{
		{ID: 1, UserID: 1, Date: time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC), Title: "Event 1"},
		{ID: 2, UserID: 1, Date: time.Date(2024, 8, 7, 0, 0, 0, 0, time.UTC), Title: "Event 2"},
		{ID: 3, UserID: 2, Date: time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC), Title: "Event 3"},
	}

	startDate := time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 8, 7, 23, 59, 59, 0, time.UTC)
	userID := 1

	filteredEvents := FilterEvents(events, userID, startDate, endDate)

	// Check the number of filtered events
	if len(filteredEvents) != 2 {
		t.Fatalf("expected 2 events, got %d", len(filteredEvents))
	}

	// Check the IDs of filtered events
	expectedIDs := []int{1, 2}
	for i, event := range filteredEvents {
		if event.ID != expectedIDs[i] {
			t.Errorf("expected event ID %d, got %d", expectedIDs[i], event.ID)
		}
	}
}
