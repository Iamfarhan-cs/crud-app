package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerCreateTaskSuccessReturnsCreated(t *testing.T) {
	recorder := performTaskRequest(
		newFakeRepository(),
		http.MethodPost,
		"/api/v1/tasks",
		`{"title":"Write handler tests"}`,
	)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}
}

func TestHandlerInvalidJSONReturnsBadRequest(t *testing.T) {
	recorder := performTaskRequest(
		newFakeRepository(),
		http.MethodPost,
		"/api/v1/tasks",
		`{"title":`,
	)

	assertErrorResponse(t, recorder, http.StatusBadRequest, "INVALID_JSON")
}

func TestHandlerInvalidTitleReturnsBadRequestAndCode(t *testing.T) {
	recorder := performTaskRequest(
		newFakeRepository(),
		http.MethodPost,
		"/api/v1/tasks",
		`{"title":"   "}`,
	)

	assertErrorResponse(t, recorder, http.StatusBadRequest, "INVALID_TITLE")
}

func TestHandlerMissingTaskReturnsNotFound(t *testing.T) {
	recorder := performTaskRequest(
		newFakeRepository(),
		http.MethodGet,
		"/api/v1/tasks/missing",
		"",
	)

	assertErrorResponse(t, recorder, http.StatusNotFound, "TASK_NOT_FOUND")
}

func TestHandlerInvalidPaginationQueryReturnsBadRequestAndCode(t *testing.T) {
	recorder := performTaskRequest(
		newFakeRepository(),
		http.MethodGet,
		"/api/v1/tasks?limit=not-an-int",
		"",
	)

	assertErrorResponse(t, recorder, http.StatusBadRequest, "INVALID_PAGINATION")
}

func TestHandlerUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	recorder := performTaskRequest(
		newFakeRepository(),
		http.MethodPut,
		"/api/v1/tasks",
		`{"title":"Nope"}`,
	)

	assertErrorResponse(t, recorder, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
}

func TestHandlerUnknownErrorMapsToInternalError(t *testing.T) {
	repo := newFakeRepository()
	repo.createErr = errors.New("database unavailable")

	recorder := performTaskRequest(
		repo,
		http.MethodPost,
		"/api/v1/tasks",
		`{"title":"Write handler tests"}`,
	)

	assertErrorResponse(t, recorder, http.StatusInternalServerError, "INTERNAL_ERROR")
}

func performTaskRequest(repo *fakeRepository, method string, target string, body string) *httptest.ResponseRecorder {
	service := NewService(repo)
	handler := NewHandler(service)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	request := httptest.NewRequest(method, target, strings.NewReader(body))
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, request)

	return recorder
}

func assertErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()

	if recorder.Code != wantStatus {
		t.Fatalf("expected status %d, got %d with body %s", wantStatus, recorder.Code, recorder.Body.String())
	}

	var body ErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if body.Error.Code != wantCode {
		t.Fatalf("expected error code %q, got %q", wantCode, body.Error.Code)
	}
}
