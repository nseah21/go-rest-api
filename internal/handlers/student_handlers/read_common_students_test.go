package student_handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"encoding/json"

	"regexp"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetCommonStudentsHandler_ReturnsHttpOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	/* Prepare mock request */
	request, _ := http.NewRequest("GET", "/api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com", nil)
	c.Request = request

	/* Prepare mock DB */
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer mockDB.Close()

	teachers := []string{
		"teacherken@gmail.com",
		"teacherjoe@gmail.com",
	}
	
	expectedRows := sqlmock.NewRows([]string{"id"}).
		AddRow("commonstudent1@gmail.com").
		AddRow("commonstudent2@gmail.com")

	sqlMock.ExpectPrepare(regexp.QuoteMeta(SELECT_COMMON_STUDENTS))
	sqlMock.ExpectQuery(regexp.QuoteMeta(SELECT_COMMON_STUDENTS)).
		WithArgs(pq.Array(teachers), len(teachers)).
		WillReturnRows(expectedRows)

	/* Prepare expected code and response */
	expectedCode := http.StatusOK
	expectedResponse := map[string][]string {
		"students": {
			"commonstudent1@gmail.com", 
			"commonstudent2@gmail.com",
		},
	}

	expectedBytes, err := json.MarshalIndent(expectedResponse, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	GetCommonStudentsHandler(mockDB)(c)

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

func TestGetCommonStudentsHandler_ReturnsHttpBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	/* Prepare mock request */
	request, _ := http.NewRequest("GET", "/api/commonstudents?cher=keyhereiswrong%40gmail.com", nil)
	c.Request = request

	/* Prepare mock DB */
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer mockDB.Close()

	sqlMock.ExpectPrepare(regexp.QuoteMeta(SELECT_COMMON_STUDENTS))

	/* Prepare expected code and response */
	expectedCode := http.StatusBadRequest
	expectedResponse := map[string]string {
		"message": GET_COMMON_STUDENTS_BAD_REQUEST_ERROR_MESSAGE,
	}
	
	expectedBytes, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Fatal(err)
	}

	/* Invoke handler */
	GetCommonStudentsHandler(mockDB)(c)

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
