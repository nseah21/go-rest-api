package notification_handlers

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"encoding/json"
	"log"
	"io"
	"bytes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRetrieveForNotificationsHandler_ReturnsHttpOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	/* Prepare mock request */
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	mockRequestBody := map[string]string {
		"teacher":  "teacherken@gmail.com",
		"notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com",
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

	expectedPrepare1 := sqlMock.ExpectPrepare(regexp.QuoteMeta(SELECT_NOTIFIABLE_STUDENTS))
	expectedPrepare1.ExpectQuery().
		WithArgs("teacherken@gmail.com", pq.Array([]string {"studentagnes@gmail.com", "studentmiche@gmail.com"})).
		WillReturnRows(
			sqlmock.NewRows([]string {"id"}).
				AddRow("studentbob@gmail.com").
				AddRow("studentagnes@gmail.com").
				AddRow("studentmiche@gmail.com"),
		)

	expectedPrepare2 := sqlMock.ExpectPrepare(regexp.QuoteMeta(INSERT_INTO_NOTIFICATIONS))
	expectedPrepare2.ExpectExec().
		WithArgs("teacherken@gmail.com", "studentbob@gmail.com", "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com").
		WillReturnResult(sqlmock.NewResult(1,1))
	expectedPrepare2.ExpectExec().
		WithArgs("teacherken@gmail.com", "studentagnes@gmail.com", "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com").
		WillReturnResult(sqlmock.NewResult(2,1))
	expectedPrepare2.ExpectExec().
		WithArgs("teacherken@gmail.com", "studentmiche@gmail.com", "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com").
		WillReturnResult(sqlmock.NewResult(3,1))

	/* Prepare expected code and response*/
	expectedCode := http.StatusOK
	expectedResponse := map[string][]string {
		"recipients": {
			"studentbob@gmail.com",
			"studentagnes@gmail.com",
			"studentmiche@gmail.com",
		},
	}

	expectedBytes, err := json.MarshalIndent(expectedResponse, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	/* Invoke handler */
	RetrieveForNotificationsHandler(mockDB)(c)

	/* Prepare actual code and response */
	actualCode := w.Code
	actualBytes, _ := io.ReadAll(w.Body)
	
	/* Check assertions */
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, string(expectedBytes), string(actualBytes))
}

func TestRetrieveForNotificationsHandler_ReturnsHttpBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	/* Prepare mock request */
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	mockRequestBody := map[string]string {
		"invalid_key":  "teacherken@gmail.com",
		"notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com",
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
	
	/* Prepare mockDB */
	sqlMock.ExpectPrepare(regexp.QuoteMeta(SELECT_NOTIFIABLE_STUDENTS))

	/* Prepare expected code and response*/
	expectedCode := http.StatusBadRequest
	expectedResponse := map[string]string {
		"message": RETRIEVE_NOTIFICATIONS_BAD_REQUEST_ERROR_MESSAGE,
	}

	expectedBytes, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Fatal(err)
	}

	/* Invoke handler */
	RetrieveForNotificationsHandler(mockDB)(c)

	/* Prepare actual code and response */
	actualCode := w.Code
	actualBytes, _ := io.ReadAll(w.Body)
	
	/* Check assertions */
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, string(expectedBytes), string(actualBytes))
}
