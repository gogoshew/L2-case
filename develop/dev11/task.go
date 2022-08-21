package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503.
	В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400.
	В случае остальных ошибок сервер должен возвращать HTTP 500.
	Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

const dateFormat = "2006-01-02"

// Logger - для логирования запросов
type Logger struct {
	handler http.Handler
}

// Конструктор логгера
func newLogger(handler http.Handler) *Logger {
	return &Logger{handler: handler}
}

// ServeHTTP логика хэндлера, опишем этот метод, чтобы удовлетворить интерфейсу
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v\n", r.Method, r.URL, time.Since(start))
}

// Event - модель JSON хранилища
type Event struct {
	UserID      int       `json:"user_id"`
	EventID     int       `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// decode декодирует данные из reader в json
func (ev *Event) decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&ev)
	if err != nil {
		log.Printf("%+v error from decoder", err)
	}
	return nil
}

// validate проверяет наличие данных в обязательных полях
func (ev *Event) validate() error {
	switch {
	case ev.UserID <= 0:
		return fmt.Errorf("invalid user_id")
	case ev.EventID <= 0:
		return fmt.Errorf("invalid event_id")
	case ev.Title == "":
		return fmt.Errorf("invalid title")
	default:
		return nil
	}
}

// Storage структура хранилища эвентов
type Storage struct {
	mu     *sync.Mutex
	events map[int][]Event
}

//Create создание события в календаре
func (s *Storage) Create(ev *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if events, ok := s.events[ev.UserID]; ok {
		for _, event := range events {
			if event.EventID == ev.EventID {
				return fmt.Errorf("%v event for %v user already exists", ev.EventID, ev.UserID)
			}
		}
	}
	s.events[ev.UserID] = append(s.events[ev.UserID], *ev)

	return nil
}

// Update обновление информации о событии в календаре
func (s *Storage) Update(ev *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var index = -1

	events := make([]Event, 0)
	ok := false

	if events, ok = s.events[ev.UserID]; !ok {
		return fmt.Errorf("user %v doesn't exist", ev.UserID)
	}

	for i, event := range events {
		if event.EventID == ev.EventID {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Errorf("can't find event with %v id for %v user id", ev.EventID, ev.UserID)
	}

	s.events[ev.UserID][index] = *ev

	return nil
}

func (s *Storage) Delete(ev *Event) (*Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := -1

	events := make([]Event, 0)
	ok := false

	if events, ok = s.events[ev.UserID]; !ok {
		return nil, fmt.Errorf("user %v doesn't exist", ev.UserID)
	}

	for i, event := range events {
		if event.EventID == ev.EventID {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Errorf("can't find event with %v id for %v user id", ev.EventID, ev.UserID)
	}

	evLen := len(s.events[ev.UserID])
	deleted := s.events[ev.UserID][index]
	s.events[ev.UserID][index] = s.events[ev.UserID][evLen-1]
	s.events[ev.UserID] = s.events[ev.UserID][:evLen-1]

	return &deleted, nil
}

func (s *Storage) getEventsForDay(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []Event

	events := make([]Event, 0)
	ok := false

	if events, ok = s.events[userID]; !ok {
		return nil, fmt.Errorf("user %v doesn't exist", userID)
	}

	for _, ev := range events {
		if ev.Date.Year() == date.Year() && ev.Date.Month() == date.Month() && ev.Date.Day() == date.Day() {
			res = append(res, ev)
		}
	}

	return res, nil
}

func (s *Storage) getEventsForWeek(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []Event

	events := make([]Event, 0)
	ok := false

	if events, ok = s.events[userID]; !ok {
		return nil, fmt.Errorf("user %v doesn't exist", userID)
	}

	for _, ev := range events {
		y1, w1 := ev.Date.ISOWeek()
		y2, w2 := date.ISOWeek()
		if y1 == y2 && w1 == w2 {
			res = append(res, ev)
		}
	}

	return res, nil
}

func (s *Storage) getEventsForMonth(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []Event

	events := make([]Event, 0)
	ok := false

	if events, ok = s.events[userID]; !ok {
		return nil, fmt.Errorf("user %v doesn't exist", userID)
	}

	for _, ev := range events {
		if ev.Date.Year() == date.Year() && ev.Date.Month() == date.Month() {
			res = append(res, ev)
		}
	}

	return res, nil
}

func getResponse(w http.ResponseWriter, r string, ev []Event, status int) {
	resp := struct {
		Result string  `json:"result"`
		Events []Event `json:"events"`
	}{Result: r, Events: ev}

	jsMarsh, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsMarsh)
}

func getErrResponse(w http.ResponseWriter, e string, status int) {
	errResp := struct {
		Error string `json:"error"`
	}{Error: e}

	jsMarsh, err := json.Marshal(errResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsMarsh)
}

// CreateEventHandler /create_event handler
func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var ev Event

	if err := ev.decode(r.Body); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ev.validate(); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.Create(&ev); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	getResponse(w, "Событие успешно создано!", []Event{ev}, http.StatusCreated)

	fmt.Println(storage.events)
}

// UpdateEventHandler /update_event handler
func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	var ev Event

	if err := ev.decode(r.Body); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ev.validate(); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.Update(&ev); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	getResponse(w, "Событие обновлено!", []Event{ev}, http.StatusOK)

	fmt.Println(storage.events)
}

// DeleteEventHandler /delete_event handler
func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	var ev Event

	if err := ev.decode(r.Body); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var deleted *Event
	var err error
	if deleted, err = storage.Delete(&ev); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
	}

	getResponse(w, "Событие удалено!", []Event{*deleted}, http.StatusOK)
}

// ForDayHandler /events_for_day handler
func ForDayHandler(w http.ResponseWriter, r *http.Request) {
	var ev []Event

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.URL.Query().Get("date"))
	if err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ev, err = storage.getEventsForDay(userID, date); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	getResponse(w, "Запрос успешно выполнен!", ev, http.StatusOK)
}

// ForWeekHandler /events_for_week handler
func ForWeekHandler(w http.ResponseWriter, r *http.Request) {
	var ev []Event

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.URL.Query().Get("date"))
	if err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ev, err = storage.getEventsForDay(userID, date); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	getResponse(w, "Запрос успешно выполнен!", ev, http.StatusOK)
}

// ForMonthHandler /events_for_month handler
func ForMonthHandler(w http.ResponseWriter, r *http.Request) {
	var ev []Event

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.URL.Query().Get("date"))
	if err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ev, err = storage.getEventsForDay(userID, date); err != nil {
		getErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	getResponse(w, "Запрос успешно выполнен!", ev, http.StatusOK)
}

// storage - глобальное хранилище событий
var storage Storage = Storage{events: make(map[int][]Event), mu: &sync.Mutex{}}

func main() {
	mux := http.NewServeMux()

	// Пропишем пути для GET
	mux.HandleFunc("/events_for_day", ForDayHandler)
	mux.HandleFunc("/events_for_week", ForWeekHandler)
	mux.HandleFunc("/events_for_month", ForMonthHandler)

	// Пропишем пути для POST
	mux.HandleFunc("/create_event", CreateEventHandler)
	mux.HandleFunc("/update_event", UpdateEventHandler)
	mux.HandleFunc("/delete_event", DeleteEventHandler)

	// Logger
	wMux := newLogger(mux)

	// Назначение порта из конфига
	port := ":8080"
	func() {
		temp := os.Getenv("PORT")
		if temp != "" {
			port = temp
		}
	}()

	log.Printf("Server is listening for requests port%v", port)
	log.Fatalln(http.ListenAndServe(port, wMux))
}
