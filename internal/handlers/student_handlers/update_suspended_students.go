package student_handlers

import (
	"database/sql"
	"encoding/json"
	"io"
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
)

func SuspendStudentHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		stmt, err := db.Prepare(UPDATE_SUSPENDED_STUDENTS)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
		}

		var suspensionRequest models.SuspensionRequest

		bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
		}

		if err := json.Unmarshal(bytes, &suspensionRequest); err != nil {
			utils.AbortWithBadRequestError(c, "Please ensure that your request body has the format { 'student': <string> }")
		}

		_, err = stmt.Exec(suspensionRequest.Student)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
