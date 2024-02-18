package student_handlers

import (
	"database/sql"
	"net/http"

	"example.com/go-rest-api/internal/models"
	"example.com/go-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	SELECT_COMMON_STUDENTS = `
		SELECT student_id 
		FROM (
			SELECT * FROM Registrations WHERE teacher_id = ANY ($1)
		) 
		GROUP BY student_id
		HAVING COUNT(teacher_id) = $2;
	`
	GET_COMMON_STUDENTS_BAD_REQUEST_ERROR_MESSAGE = `Please specify the parameters using the key 'teacher'`
)

func GetCommonStudentsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		stmt, err := db.Prepare(SELECT_COMMON_STUDENTS)

		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
			return
		}

		var students []string

		parameters := c.Request.URL.Query()
		if teachers, found := parameters["teacher"]; found {
			rows, err := stmt.Query(pq.Array(teachers), len(teachers))
			if err != nil {
				utils.AbortWithInternalServerError(c, err.Error())
				return
			}
			defer rows.Close()

			for rows.Next() {
				var student models.Student
				if err := rows.Scan(&student.Id); err != nil {
					utils.AbortWithInternalServerError(c, err.Error())
				}
				students = append(students, student.Id)
			}
		} else {
			utils.AbortWithBadRequestError(c, GET_COMMON_STUDENTS_BAD_REQUEST_ERROR_MESSAGE)
			return
		}

		result := map[string][]string{
			"students": students,
		}

		c.IndentedJSON(http.StatusOK, result)
	}
}
