package registration_handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRegisterStudentsHandler_ReturnsHttpNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	/* Prepare mock request */
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	mockRequestBody := map[string]any {
		"teacher": "teacher1@gmail.com",
		"students": []string {"student1@gmail.com"},
	}

	jsonBytes, err := json.Marshal(mockRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	/* Prepare mock DB */
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer mockDB.Close()

	expectedPrepare := sqlMock.ExpectPrepare(regexp.QuoteMeta(INSERT_INTO_REGISTRATIONS))
	expectedPrepare.ExpectExec().
		WithArgs("teacher1@gmail.com", "student1@gmail.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	/* Prepare expected code and response */
	expectedCode := http.StatusNoContent
	expectedBytes := make([]byte, 0)

	/* Invoke handler */
	RegisterStudentsHandler(mockDB)(c)

	/* Prepare actual code and response */
	actualCode := w.Code
	actualBytes, _ := io.ReadAll(w.Body)

	/* Check assertions */
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedBytes, actualBytes)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestRegisterStudentsHandler_ReturnsHttpBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	/* Prepare mock request */
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	mockRequestBody := map[string]string {
		"teacher": "teacher1@gmail.com",
		"students": "student1@gmail.com",
	}

	jsonBytes, err := json.Marshal(mockRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	/* Prepare mockDB */
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer mockDB.Close()

	sqlMock.ExpectPrepare(regexp.QuoteMeta(INSERT_INTO_REGISTRATIONS))

	/* Prepare expected code and response */
	expectedCode := http.StatusBadRequest
	expectedResponse := map[string]string {
		"message": "Please ensure that your JSON matches the following format: { teacher: string, students: string[] }",
	}

	expectedBytes, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Fatal(err)
	}

	/* Invoke handler */
	RegisterStudentsHandler(mockDB)(c)

	/* Prepare actual code and response */
	actualCode := w.Code
	actualBytes, _ := io.ReadAll(w.Body)

	/* Check assertions */
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, string(expectedBytes), string(actualBytes))

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
