package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestBusinessError_Error(t *testing.T) {
	err := &BusinessError{
		Code:    1001,
		Message: "test error",
	}

	expected := "code: 1001, message: test error"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNew(t *testing.T) {
	err := New(1001, "test message")

	if err.Code != 1001 {
		t.Errorf("expected code 1001, got %d", err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("expected message 'test message', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusOK {
		t.Errorf("expected HTTPStatus 200, got %d", err.HTTPStatus)
	}
}

func TestNewWithHTTPStatus(t *testing.T) {
	err := NewWithHTTPStatus(401, "unauthorized", http.StatusUnauthorized)

	if err.Code != 401 {
		t.Errorf("expected code 401, got %d", err.Code)
	}
	if err.Message != "unauthorized" {
		t.Errorf("expected message 'unauthorized', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusUnauthorized {
		t.Errorf("expected HTTPStatus 401, got %d", err.HTTPStatus)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		message        string
		expectedCode   int
		expectedPrefix string
	}{
		{
			name:           "wrap business error",
			err:            New(1001, "original"),
			message:        "wrapped",
			expectedCode:   1001,
			expectedPrefix: "wrapped: original",
		},
		{
			name:           "wrap standard error",
			err:            errors.New("standard error"),
			message:        "wrapped",
			expectedCode:   CodeInternalError,
			expectedPrefix: "wrapped: standard error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := Wrap(tt.err, tt.message)
			if wrapped.Code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, wrapped.Code)
			}
			if wrapped.Message != tt.expectedPrefix {
				t.Errorf("expected message '%s', got '%s'", tt.expectedPrefix, wrapped.Message)
			}
		})
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   *BusinessError
		expected bool
	}{
		{
			name:     "matching business error",
			err:      New(1001, "error"),
			target:   New(1001, "different message"),
			expected: true,
		},
		{
			name:     "non-matching business error",
			err:      New(1001, "error"),
			target:   New(1002, "error"),
			expected: false,
		},
		{
			name:     "standard error",
			err:      errors.New("standard"),
			target:   New(1001, "error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Is(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		name               string
		fn                 func(string) *BusinessError
		expectedCode       int
		expectedHTTPStatus int
	}{
		{"BadRequest", BadRequest, CodeBadRequest, http.StatusBadRequest},
		{"Unauthorized", Unauthorized, CodeUnauthorized, http.StatusUnauthorized},
		{"Forbidden", Forbidden, CodeForbidden, http.StatusForbidden},
		{"NotFound", NotFound, CodeNotFound, http.StatusNotFound},
		{"InternalError", InternalError, CodeInternalError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn("test message")
			if err.Code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, err.Code)
			}
			if err.HTTPStatus != tt.expectedHTTPStatus {
				t.Errorf("expected HTTPStatus %d, got %d", tt.expectedHTTPStatus, err.HTTPStatus)
			}
			if err.Message != "test message" {
				t.Errorf("expected message 'test message', got '%s'", err.Message)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name               string
		err                *BusinessError
		expectedCode       int
		expectedHTTPStatus int
	}{
		{"ErrBadRequest", ErrBadRequest, CodeBadRequest, http.StatusOK},
		{"ErrUnauthorized", ErrUnauthorized, CodeUnauthorized, http.StatusUnauthorized},
		{"ErrForbidden", ErrForbidden, CodeForbidden, http.StatusForbidden},
		{"ErrNotFound", ErrNotFound, CodeNotFound, http.StatusOK},
		{"ErrInternalError", ErrInternalError, CodeInternalError, http.StatusOK},
		{"ErrParamError", ErrParamError, CodeParamError, http.StatusOK},
		{"ErrDataNotFound", ErrDataNotFound, CodeDataNotFound, http.StatusOK},
		{"ErrDataExists", ErrDataExists, CodeDataExists, http.StatusOK},
		{"ErrOperationFailed", ErrOperationFailed, CodeOperationFailed, http.StatusOK},
		{"ErrPermissionDenied", ErrPermissionDenied, CodePermissionDenied, http.StatusOK},
		{"ErrTokenExpired", ErrTokenExpired, CodeTokenExpired, http.StatusUnauthorized},
		{"ErrTokenInvalid", ErrTokenInvalid, CodeTokenInvalid, http.StatusUnauthorized},
		{"ErrLoginFailed", ErrLoginFailed, CodeLoginFailed, http.StatusOK},
		{"ErrCaptchaError", ErrCaptchaError, CodeCaptchaError, http.StatusOK},
		{"ErrAccountDisabled", ErrAccountDisabled, CodeAccountDisabled, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, tt.err.Code)
			}
			if tt.err.HTTPStatus != tt.expectedHTTPStatus {
				t.Errorf("expected HTTPStatus %d, got %d", tt.expectedHTTPStatus, tt.err.HTTPStatus)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Verify error codes are correct
	if CodeSuccess != 200 {
		t.Errorf("expected CodeSuccess 200, got %d", CodeSuccess)
	}
	if CodeBadRequest != 400 {
		t.Errorf("expected CodeBadRequest 400, got %d", CodeBadRequest)
	}
	if CodeUnauthorized != 401 {
		t.Errorf("expected CodeUnauthorized 401, got %d", CodeUnauthorized)
	}
	if CodeForbidden != 403 {
		t.Errorf("expected CodeForbidden 403, got %d", CodeForbidden)
	}
	if CodeNotFound != 404 {
		t.Errorf("expected CodeNotFound 404, got %d", CodeNotFound)
	}
	if CodeInternalError != 500 {
		t.Errorf("expected CodeInternalError 500, got %d", CodeInternalError)
	}

	// Business error codes
	if CodeParamError != 1001 {
		t.Errorf("expected CodeParamError 1001, got %d", CodeParamError)
	}
	if CodeTokenExpired != 2001 {
		t.Errorf("expected CodeTokenExpired 2001, got %d", CodeTokenExpired)
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(1001, "benchmark error")
	}
}

func BenchmarkWrap(b *testing.B) {
	err := New(1001, "original error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Wrap(err, "wrapped")
	}
}

func BenchmarkIs(b *testing.B) {
	err := New(1001, "test error")
	target := New(1001, "target error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Is(err, target)
	}
}
