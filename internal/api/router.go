package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type ContactInput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
}

type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
}

type WebhookInput struct {
	Provider  string          `json:"provider"`
	EventID   string          `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

type Event struct {
	Provider  string          `json:"provider"`
	EventID   string          `json:"event_id"`
	EventType string          `json:"event_type"`
	Received  time.Time       `json:"received_at"`
	Payload   json.RawMessage `json:"payload"`
}

type Store struct {
	mu     sync.RWMutex
	events map[string]Event
}

func NewStore() *Store { return &Store{events: make(map[string]Event)} }

func (s *Store) Add(input WebhookInput) (bool, Event) {
	key := strings.ToLower(strings.TrimSpace(input.Provider)) + ":" + strings.TrimSpace(input.EventID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if event, exists := s.events[key]; exists {
		return false, event
	}
	event := Event{
		Provider: strings.ToLower(strings.TrimSpace(input.Provider)),
		EventID: strings.TrimSpace(input.EventID),
		EventType: strings.TrimSpace(input.EventType),
		Received: time.Now().UTC(),
		Payload: input.Payload,
	}
	s.events[key] = event
	return true, event
}

func (s *Store) List(provider string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Event, 0, len(s.events))
	for _, event := range s.events {
		if provider == "" || event.Provider == provider {
			out = append(out, event)
		}
	}
	return out
}

func NewRouter() http.Handler {
	store := NewStore()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "relayops-api"})
	})
	mux.HandleFunc("POST /api/v1/contacts/normalize", normalizeContact)
	mux.HandleFunc("POST /api/v1/webhooks/ingest", ingestWebhook(store))
	mux.HandleFunc("GET /api/v1/events", listEvents(store))
	return mux
}

func normalizeContact(w http.ResponseWriter, r *http.Request) {
	var input ContactInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON.")
		return
	}
	contact, err := normalize(input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "invalid_contact", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, contact)
}

func normalize(input ContactInput) (Contact, error) {
	contact := Contact{
		Name: strings.Join(strings.Fields(input.Name), " "),
		Email: strings.ToLower(strings.TrimSpace(input.Email)),
		Company: strings.Join(strings.Fields(input.Company), " "),
	}
	if contact.Name == "" {
		return Contact{}, errors.New("name is required")
	}
	if !emailPattern.MatchString(contact.Email) {
		return Contact{}, errors.New("a valid email is required")
	}
	return contact, nil
}

func ingestWebhook(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input WebhookInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON.")
			return
		}
		input.Provider = strings.TrimSpace(input.Provider)
		input.EventID = strings.TrimSpace(input.EventID)
		input.EventType = strings.TrimSpace(input.EventType)
		if input.Provider == "" || input.EventID == "" || input.EventType == "" || len(input.Payload) == 0 || !json.Valid(input.Payload) {
			writeError(w, http.StatusUnprocessableEntity, "invalid_event", "provider, event_id, event_type, and a valid JSON payload are required.")
			return
		}
		accepted, event := store.Add(input)
		status := http.StatusAccepted
		if !accepted {
			status = http.StatusOK
		}
		writeJSON(w, status, map[string]any{"accepted": accepted, "duplicate": !accepted, "event": event})
	}
}

func listEvents(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("provider")))
		writeJSON(w, http.StatusOK, map[string]any{"events": store.List(provider)})
	}
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"error": code, "message": message})
}
