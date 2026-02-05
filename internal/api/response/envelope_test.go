package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func newTestContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestOK(t *testing.T) {
	c, rec := newTestContext()
	err := OK(c, map[string]string{"key": "value"})
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var env Envelope
	json.Unmarshal(rec.Body.Bytes(), &env)
	if env.Code != 0 {
		t.Errorf("expected code 0, got %d", env.Code)
	}
	if env.Message != "success" {
		t.Errorf("expected message 'success', got '%s'", env.Message)
	}
	if env.Timestamp == 0 {
		t.Error("timestamp should be set")
	}
}

func TestCreated(t *testing.T) {
	c, rec := newTestContext()
	err := Created(c, "item")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestPaged(t *testing.T) {
	c, rec := newTestContext()
	items := []string{"a", "b"}
	err := Paged(c, items, 50, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	var env Envelope
	json.Unmarshal(rec.Body.Bytes(), &env)
	dataMap := env.Data.(map[string]interface{})
	if int64(dataMap["total"].(float64)) != 50 {
		t.Error("expected total 50")
	}
}

func TestNoContent(t *testing.T) {
	c, rec := newTestContext()
	err := NoContent(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestErr(t *testing.T) {
	c, rec := newTestContext()
	err := Err(c, http.StatusBadRequest, 40001, "invalid input", map[string]string{"field": "name"})
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	var env Envelope
	json.Unmarshal(rec.Body.Bytes(), &env)
	if env.Code != 40001 {
		t.Errorf("expected code 40001, got %d", env.Code)
	}
}
