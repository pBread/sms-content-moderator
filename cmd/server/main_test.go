package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pBread/sms-content-moderator/internal/evaluator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the evaluator to control its behavior during tests.
type MockEvaluator struct {
	mock.Mock
}

func (m *MockEvaluator) EvaluateContent(content string) (evaluator.Response, error) {
	args := m.Called(content)
	return args.Get(0).(evaluator.Response), args.Error(1)
}

func TestUnauthenticatedHandler(t *testing.T) {
	// Initialize the mock and replace the actual evaluator with the mock.
	mockEval := new(MockEvaluator)
	eval = mockEval

	tests := []struct {
		name           string
		method         string
		body           string
		mockResponse   interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Reject non-POST methods",
			method:         "GET",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Only POST method is allowed\n",
		},
		{
			name:           "Bad request on invalid JSON",
			method:         "POST",
			body:           "{invalid-json}",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error reading request body\n",
		},
		{
			name:           "Internal server error on evaluator error",
			method:         "POST",
			body:           `{"Message":"test message"}`,
			mockError:      errors.New("internal error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "internal error\n",
		},
		{
			name:           "Successful request",
			method:         "POST",
			body:           `{"Message":"test message"}`,
			mockResponse:   map[string]interface{}{"Result": "clean"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Result":"clean"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup request and recorder
			req, err := http.NewRequest(tc.method, "/evaluate-message", bytes.NewBufferString(tc.body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			// Setup handler
			handler := http.HandlerFunc(unauthenticatedHandler)

			// Set expected responses from the mock
			mockEval.On("EvaluateContent", "test message").Return(tc.mockResponse, tc.mockError)

			// Serve HTTP
			handler.ServeHTTP(rr, req)

			// Assert HTTP status
			assert.Equal(t, tc.expectedStatus, rr.Code)
			// Assert HTTP body
			assert.JSONEq(t, tc.expectedBody, rr.Body.String())

			// Clear expectations on the mock
			mockEval.ExpectedCalls = nil
			mockEval.Calls = nil
		})
	}
}
