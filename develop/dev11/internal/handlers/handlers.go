package handlers

import (
	"net/http"
	"strconv"

	"dev11/internal/db"
	"dev11/internal/utils"
)

var currentID int

// создание события
func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	event, err := utils.ParseFormAndValidate(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	currentID++
	event.ID = currentID
	db.Events = append(db.Events, event)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"result": "event created"})
}

// обновление события
func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	event, err := utils.ParseFormAndValidate(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	for i, e := range db.Events {
		if e.ID == id {
			db.Events[i] = event
			db.Events[i].ID = id
			utils.WriteJSON(w, http.StatusOK, map[string]string{"result": "event updated"})
			return
		}
	}

	utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "event not found"})
}

// удаление события
func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	for i, e := range db.Events {
		if e.ID == id {
			db.Events = append(db.Events[:i], db.Events[i+1:]...)
			utils.WriteJSON(w, http.StatusOK, map[string]string{"result": "event deleted"})
			return
		}
	}

	utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "event not found"})
}

// получение событий за день
func EventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := utils.ParseGetParams(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	filteredEvents := utils.FilterEvents(db.Events, userID, date, date)
	utils.WriteJSON(w, http.StatusOK, filteredEvents)
}

// получение событий за неделю
func EventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := utils.ParseGetParams(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	endDate := date.AddDate(0, 0, 7)
	filteredEvents := utils.FilterEvents(db.Events, userID, date, endDate)
	utils.WriteJSON(w, http.StatusOK, filteredEvents)
}

// получение событий за месяц
func EventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := utils.ParseGetParams(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	endDate := date.AddDate(0, 1, 0)
	filteredEvents := utils.FilterEvents(db.Events, userID, date, endDate)
	utils.WriteJSON(w, http.StatusOK, filteredEvents)
}
