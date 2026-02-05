package apperr

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(CodeInvalidInput, "bad request")
	if err.Code != CodeInvalidInput {
		t.Errorf("expected code %d, got %d", CodeInvalidInput, err.Code)
	}
	if err.Message != "bad request" {
		t.Errorf("expected message 'bad request', got '%s'", err.Message)
	}
	if err.Cause != nil {
		t.Error("expected nil cause")
	}
}

func TestWrap(t *testing.T) {
	cause := fmt.Errorf("db connection lost")
	err := Wrap(cause, CodeDBError, "database error")
	if err.Code != CodeDBError {
		t.Errorf("expected code %d, got %d", CodeDBError, err.Code)
	}
	if !errors.Is(err, cause) {
		t.Error("expected Unwrap to return cause")
	}
	if err.Error() != fmt.Sprintf("[%d] database error: db connection lost", CodeDBError) {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("user", "abc-123")
	if err.Code != CodeNotFound {
		t.Errorf("expected code %d, got %d", CodeNotFound, err.Code)
	}
	if !IsNotFound(err) {
		t.Error("IsNotFound should return true")
	}
}

func TestConflict(t *testing.T) {
	err := Conflict("already exists")
	if !IsConflict(err) {
		t.Error("IsConflict should return true")
	}
}

func TestHasRelation(t *testing.T) {
	err := HasRelation("has children")
	if !IsConflict(err) {
		t.Error("IsConflict should return true for HasRelation")
	}
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("token expired")
	if !IsUnauthorized(err) {
		t.Error("IsUnauthorized should return true")
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("no permission")
	if !IsForbidden(err) {
		t.Error("IsForbidden should return true")
	}
}

func TestWithDetails(t *testing.T) {
	err := InvalidInput("validation failed")
	detailed := WithDetails(err, map[string]interface{}{"field": "email"})
	if detailed.Details["field"] != "email" {
		t.Error("expected details to contain field=email")
	}
	if err.Details != nil {
		t.Error("original error should not be modified")
	}
}

func TestInternal(t *testing.T) {
	cause := fmt.Errorf("panic")
	err := Internal("something went wrong", cause)
	if err.Code != CodeInternal {
		t.Errorf("expected code %d, got %d", CodeInternal, err.Code)
	}
	if !errors.Is(err, cause) {
		t.Error("should unwrap to cause")
	}
}

func TestErrorWithNilCause(t *testing.T) {
	err := New(CodeInvalidInput, "test")
	expected := fmt.Sprintf("[%d] test", CodeInvalidInput)
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestIsChecksWithNonAppError(t *testing.T) {
	err := fmt.Errorf("plain error")
	if IsNotFound(err) {
		t.Error("should return false for non-apperr")
	}
	if IsConflict(err) {
		t.Error("should return false for non-apperr")
	}
	if IsUnauthorized(err) {
		t.Error("should return false for non-apperr")
	}
	if IsForbidden(err) {
		t.Error("should return false for non-apperr")
	}
}
