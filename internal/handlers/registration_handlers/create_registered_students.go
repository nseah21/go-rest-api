package registration_handlers

import (
	"database/sql"
	"net/http"

	"example.com/go-rest-api/internal/models"
	"example.com/go-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	INSERT_INTO_REGISTRATIONS = `
		INSERT INTO Registrations (teacher_id, student_id)
		VALUES ($1, $2);
	`
	REGISTER_STUDENT_BAD_REQUEST_ERROR_MESSAGE = `Please ensure that your JSON matches the following format: { teacher: string, students: string[] }`
)

func RegisterStudentsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		stmt, err := db.Prepare(INSERT_INTO_REGISTRATIONS)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
			return
		}

		var registration models.RegistrationRequest

		if err := c.BindJSON(&registration); err != nil {
			utils.AbortWithBadRequestError(c, REGISTER_STUDENT_BAD_REQUEST_ERROR_MESSAGE)
			return
		}

		for _, student := range registration.Students {
			_, err = stmt.Exec(registration.Teacher, student)
			if err != nil {
				utils.AbortWithInternalServerError(c, err.Error())
				return
			}
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
