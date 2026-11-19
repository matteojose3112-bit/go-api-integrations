package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func request(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestHealth(t *testing.T) {
	rec := request(t, NewRouter(), http.MethodGet, "/health", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestNormalizeContact(t *testing.T) {
	rec := request(t, NewRouter(), http.MethodPost, "/api/v1/contacts/normalize", `{"name":"  Ada   Lovelace ","email":" ADA@EXAMPLE.COM ","company":"  Analytical   Engines "}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var contact Contact
	if err := json.Unmarshal(rec.Body.Bytes(), &contact); err != nil {
		t.Fatal(err)
	}
	if contact.Name != "Ada Lovelace" || contact.Email != "ada@example.com" || contact.Company != "Analytical Engines" {
		t.Fatalf("unexpected normalized contact: %#v", contact)
	}
}

func TestRejectInvalidContact(t *testing.T) {
	rec := request(t, NewRouter(), http.MethodPost, "/api/v1/contacts/normalize", `{"name":"Ada","email":"not-an-email"}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}

func TestWebhookIdempotency(t *testing.T) {
	handler := NewRouter()
	body := `{"provider":"CRM","event_id":"evt-123","event_type":"contact.created","payload":{"id":7}}`
	first := request(t, handler, http.MethodPost, "/api/v1/webhooks/ingest", body)
	second := request(t, handler, http.MethodPost, "/api/v1/webhooks/ingest", body)
	if first.Code != http.StatusAccepted || second.Code != http.StatusOK {
		t.Fatalf("unexpected statuses: %d and %d", first.Code, second.Code)
	}
	if !bytes.Contains(second.Body.Bytes(), []byte(`"duplicate":true`)) {
		t.Fatalf("expected duplicate marker, got %s", second.Body.String())
	}
}

func TestListEvents(t *testing.T) {
	handler := NewRouter()
	request(t, handler, http.MethodPost, "/api/v1/webhooks/ingest", `{"provider":"crm","event_id":"evt-1","event_type":"updated","payload":{"ok":true}}`)
	rec := request(t, handler, http.MethodGet, "/api/v1/events?provider=crm", "")
	if rec.Code != http.StatusOK || !bytes.Contains(rec.Body.Bytes(), []byte("evt-1")) {
		t.Fatalf("unexpected event list: %d %s", rec.Code, rec.Body.String())
	}
}
