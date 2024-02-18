package student_handlers

import (
	"database/sql"
	"net/http"

	"example.com/go-rest-api/internal/models"
	"example.com/go-rest-api/internal/utils"
	"github.com/gin-gonic/gin"	
	_ "github.com/lib/pq"
)

const (
	UPDATE_SUSPENDED_STUDENTS = `
		UPDATE Students
		SET has_been_suspended = true
		WHERE id = $1;
	`
	SUSPEND_STUDENTS_BAD_REQUEST_ERROR_MESSAGE = `Please ensure that your request body has the format { 'student': <string> }`
)

func SuspendStudentHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		stmt, err := db.Prepare(UPDATE_SUSPENDED_STUDENTS)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
			return
		}

		var suspensionRequest models.SuspensionRequest

		if err := c.BindJSON(&suspensionRequest); err != nil {
			utils.AbortWithBadRequestError(c, SUSPEND_STUDENTS_BAD_REQUEST_ERROR_MESSAGE)
			return
		}

		_, err = stmt.Exec(suspensionRequest.Student)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
