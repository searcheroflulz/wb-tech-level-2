package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

type Event struct {
	UserID      string `json:"user_id"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

func parseEventParams(r *http.Request) (*Event, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	_, err = strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return nil, err
	}

	date := r.FormValue("date")
	if _, err = time.Parse("2006-01-02", date); err != nil {
		return nil, err
	}
	description := r.FormValue("description")

	return &Event{
		UserID:      r.FormValue("user_id"),
		Date:        date,
		Description: description,
	}, nil
}

type Calendar struct {
	Events map[string][]Event
	sync.RWMutex
}

type APIResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

var cache = Calendar{Events: make(map[string][]Event)}

func (c *Calendar) AddEvent(data Event) {
	date := data.Date
	c.Lock()
	defer c.Unlock()
	c.Events[date] = append(c.Events[date], data)
}

func (c *Calendar) DeleteEvent(data Event) bool {
	var updatedEvents []Event
	var lenEvents int

	c.RLock()

	if events, ok := c.Events[data.Date]; ok {
		lenEvents = len(events)
		for _, event := range events {
			if event.UserID != data.UserID {
				updatedEvents = append(updatedEvents, event)
			}
		}
	}
	c.RUnlock()

	if lenEvents != len(updatedEvents) {
		c.Lock()
		c.Events[data.Date] = updatedEvents
		c.Unlock()
		return true
	}

	return false
}

func (c *Calendar) UpdateEvent(data Event) bool {
	if events, ok := c.Events[data.Date]; ok {
		for i, event := range events {
			if event.UserID == data.UserID {
				c.Events[data.Date][i] = data
				return true
			}
		}
	}

	return false
}

func (c *Calendar) GetEventsDay(data Event) []Event {
	var result []Event

	c.RLock()
	defer c.RUnlock()

	for _, event := range c.Events[data.Date] {
		if event.UserID == data.UserID {
			result = append(result, event)
		}
	}

	return result
}

func (c *Calendar) GetEventsWeek(data Event) []Event {
	eventDate, _ := time.Parse("2006-01-02", data.Date)
	weekday := eventDate.Weekday()
	startWeek := eventDate.AddDate(0, 0, -int(weekday))

	endWeek := startWeek.AddDate(0, 0, 7)
	var result []Event

	c.RLock()
	defer c.RUnlock()

	for date, events := range c.Events {
		eventDate, _ := time.Parse("2006-01-02", date)

		if eventDate.After(startWeek) && eventDate.Before(endWeek) {
			for _, event := range events {
				if event.UserID == data.UserID {
					result = append(result, event)
				}
			}
		}
	}

	return result
}

func (c *Calendar) GetEventsMonth(data Event) []Event {
	eventDate, _ := time.Parse("2006-01-02", data.Date)

	startMonth := time.Date(eventDate.Year(), eventDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := startMonth.AddDate(0, 1, 0)
	var result []Event

	c.RLock()
	defer c.RUnlock()

	for date, events := range c.Events {
		eventDate, _ := time.Parse("2006-01-02", date)

		if eventDate.After(startMonth) && eventDate.Before(endMonth) {
			for _, event := range events {
				if event.UserID == data.UserID {
					result = append(result, event)
				}
			}
		}
	}

	return result
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := APIResponse{Error: "Wrong method"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	event, err := parseEventParams(r)
	if err != nil {
		response := APIResponse{Error: "Invalid parameters"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	cache.AddEvent(*event)

	response := APIResponse{Result: "Successfully added event"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := APIResponse{Error: "Wrong method"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	event, err := parseEventParams(r)
	if err != nil {
		response := APIResponse{Error: "Invalid parameters"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if cache.UpdateEvent(*event) {
		response := APIResponse{Result: "Event updated successfully"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{Error: "Event not found"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := APIResponse{Error: "Wrong method"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	event, err := parseEventParams(r)
	if err != nil {
		response := APIResponse{Error: "Invalid parameters"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if cache.DeleteEvent(*event) {
		response := APIResponse{Result: "Event deleted successfully"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{Error: "Event not found"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}

func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := APIResponse{Error: "Wrong method"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	event, err := parseEventParams(r)
	if err != nil {
		response := APIResponse{Error: "Invalid parameters"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	events := cache.GetEventsWeek(*event)

	if len(events) == 0 {
		response := APIResponse{Error: "No events found for the specified user and date"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{Result: events}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := APIResponse{Error: "Wrong method"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	event, err := parseEventParams(r)
	if err != nil {
		response := APIResponse{Error: "Invalid parameters"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	events := cache.GetEventsDay(*event)

	if len(events) == 0 {
		response := APIResponse{Error: "No events found for the specified user and date"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{Result: events}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := APIResponse{Error: "Wrong method"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	event, err := parseEventParams(r)
	if err != nil {
		response := APIResponse{Error: "Invalid parameters"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	events := cache.GetEventsMonth(*event)

	if len(events) == 0 {
		response := APIResponse{Error: "No events found for the specified user and date"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{Result: events}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func startServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", logger(createEventHandler))
	mux.HandleFunc("/update_event", logger(updateEventHandler))
	mux.HandleFunc("/delete_event", logger(deleteEventHandler))

	mux.HandleFunc("/events_for_day", logger(eventsForDayHandler))
	mux.HandleFunc("/events_for_week", logger(eventsForWeekHandler))
	mux.HandleFunc("/events_for_month", logger(eventsForMonthHandler))

	log.Printf("Starting server on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s\n", r.Method, r.RequestURI, r.RemoteAddr)
		next(w, r)
	}
}

func main() {
	port := flag.Int("port", 8080, "Port for the server")
	flag.Parse()

	startServer(fmt.Sprintf(":%d", *port))
}
