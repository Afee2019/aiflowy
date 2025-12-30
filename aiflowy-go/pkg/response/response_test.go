package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// Helper function to create test Echo context
func createTestContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// Helper function to parse response body
func parseResponse(rec *httptest.ResponseRecorder) (*Response, error) {
	var resp Response
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	return &resp, err
}

func TestSuccess(t *testing.T) {
	c, rec := createTestContext()

	data := map[string]string{"key": "value"}
	err := Success(c, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	resp, err := parseResponse(rec)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("expected message 'success', got '%s'", resp.Message)
	}
}

func TestSuccessWithMessage(t *testing.T) {
	c, rec := createTestContext()

	err := SuccessWithMessage(c, "custom message", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := parseResponse(rec)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
	if resp.Message != "custom message" {
		t.Errorf("expected message 'custom message', got '%s'", resp.Message)
	}
}

func TestCreated(t *testing.T) {
	c, rec := createTestContext()

	data := map[string]string{"id": "123"}
	err := Created(c, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}

	resp, err := parseResponse(rec)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
}

func TestNoContent(t *testing.T) {
	c, rec := createTestContext()

	err := NoContent(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}
}

func TestPageSuccess(t *testing.T) {
	c, rec := createTestContext()

	items := []map[string]string{{"id": "1"}, {"id": "2"}}
	err := PageSuccess(c, items, 100, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	resp, err := parseResponse(rec)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}

	// Parse data as PageData
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected data to be map, got %T", resp.Data)
	}
	if dataMap["total"].(float64) != 100 {
		t.Errorf("expected total 100, got %v", dataMap["total"])
	}
	if dataMap["page"].(float64) != 1 {
		t.Errorf("expected page 1, got %v", dataMap["page"])
	}
	if dataMap["size"].(float64) != 10 {
		t.Errorf("expected size 10, got %v", dataMap["size"])
	}
}

func TestError(t *testing.T) {
	c, rec := createTestContext()

	err := Error(c, http.StatusInternalServerError, "internal error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected code 500, got %d", resp.Code)
	}
	if resp.Message != "internal error" {
		t.Errorf("expected message 'internal error', got '%s'", resp.Message)
	}
}

func TestBadRequest(t *testing.T) {
	c, rec := createTestContext()

	err := BadRequest(c, "bad request message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected code 400, got %d", resp.Code)
	}
}

func TestUnauthorized(t *testing.T) {
	c, rec := createTestContext()

	err := Unauthorized(c, "unauthorized message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected code 401, got %d", resp.Code)
	}
}

func TestForbidden(t *testing.T) {
	c, rec := createTestContext()

	err := Forbidden(c, "forbidden message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusForbidden {
		t.Errorf("expected code 403, got %d", resp.Code)
	}
}

func TestNotFound(t *testing.T) {
	c, rec := createTestContext()

	err := NotFound(c, "not found message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected code 404, got %d", resp.Code)
	}
}

func TestConflict(t *testing.T) {
	c, rec := createTestContext()

	err := Conflict(c, "conflict message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusConflict {
		t.Errorf("expected code 409, got %d", resp.Code)
	}
}

func TestUnprocessableEntity(t *testing.T) {
	c, rec := createTestContext()

	validationErrors := []ValidationError{
		{Field: "name", Message: "name is required"},
		{Field: "email", Message: "invalid email format"},
	}
	err := UnprocessableEntity(c, validationErrors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected code 422, got %d", resp.Code)
	}
}

func TestTooManyRequests(t *testing.T) {
	c, rec := createTestContext()

	err := TooManyRequests(c, "rate limit exceeded")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusTooManyRequests {
		t.Errorf("expected code 429, got %d", resp.Code)
	}
}

func TestInternalError(t *testing.T) {
	c, rec := createTestContext()

	err := InternalError(c, "server error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected code 500, got %d", resp.Code)
	}
}

func TestServiceUnavailable(t *testing.T) {
	c, rec := createTestContext()

	err := ServiceUnavailable(c, "service down")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("expected code 503, got %d", resp.Code)
	}
}

func TestOK(t *testing.T) {
	c, rec := createTestContext()

	data := map[string]string{"key": "value"}
	err := OK(c, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
}

func TestOKWithMessage(t *testing.T) {
	c, rec := createTestContext()

	err := OKWithMessage(c, "custom", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, _ := parseResponse(rec)
	if resp.Message != "custom" {
		t.Errorf("expected message 'custom', got '%s'", resp.Message)
	}
}

func TestFail(t *testing.T) {
	c, rec := createTestContext()

	err := Fail(c, 1001, "business error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Fail always returns HTTP 200
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	resp, _ := parseResponse(rec)
	if resp.Code != 1001 {
		t.Errorf("expected code 1001, got %d", resp.Code)
	}
	if resp.Message != "business error" {
		t.Errorf("expected message 'business error', got '%s'", resp.Message)
	}
}

func TestResponseStructure(t *testing.T) {
	// Test Response struct JSON marshaling
	resp := Response{
		Code:    0,
		Message: "success",
		Data:    map[string]string{"key": "value"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	var parsed Response
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if parsed.Code != 0 || parsed.Message != "success" {
		t.Error("response marshaling/unmarshaling failed")
	}
}

func TestPageDataStructure(t *testing.T) {
	// Test PageData struct JSON marshaling
	pageData := PageData{
		Items: []string{"item1", "item2"},
		Total: 100,
		Page:  1,
		Size:  10,
	}

	data, err := json.Marshal(pageData)
	if err != nil {
		t.Fatalf("failed to marshal page data: %v", err)
	}

	var parsed PageData
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal page data: %v", err)
	}

	if parsed.Total != 100 || parsed.Page != 1 || parsed.Size != 10 {
		t.Error("page data marshaling/unmarshaling failed")
	}
}

func TestValidationErrorStructure(t *testing.T) {
	// Test ValidationError struct JSON marshaling
	validationErr := ValidationError{
		Field:   "email",
		Message: "invalid format",
		Code:    "INVALID_FORMAT",
	}

	data, err := json.Marshal(validationErr)
	if err != nil {
		t.Fatalf("failed to marshal validation error: %v", err)
	}

	var parsed ValidationError
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal validation error: %v", err)
	}

	if parsed.Field != "email" || parsed.Message != "invalid format" || parsed.Code != "INVALID_FORMAT" {
		t.Error("validation error marshaling/unmarshaling failed")
	}
}

// Benchmarks

func BenchmarkSuccess(b *testing.B) {
	e := echo.New()
	data := map[string]string{"key": "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = Success(c, data)
	}
}

func BenchmarkPageSuccess(b *testing.B) {
	e := echo.New()
	items := []map[string]string{{"id": "1"}, {"id": "2"}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = PageSuccess(c, items, 100, 1, 10)
	}
}

func BenchmarkError(b *testing.B) {
	e := echo.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = Error(c, http.StatusBadRequest, "error message")
	}
}
