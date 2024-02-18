package notification_handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"example.com/go-rest-api/internal/models"
	"example.com/go-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	SELECT_NOTIFIABLE_STUDENTS = `
		SELECT S.id 
		FROM Students S JOIN Registrations R on S.id = R.student_id 
		WHERE S.has_been_suspended = false and R.teacher_id = $1; 
	`
	INSERT_INTO_NOTIFICATIONS = `
		INSERT INTO Notifications (sender_id, recipient_id, notification) 
		VALUES ($1, $2, $3);
	`
	RETRIEVE_NOTIFICATIONS_BAD_REQUEST_ERROR_MESSAGE = `Please ensure that your request body has the following format: { 'teacher': <string>, 'notification': <string> }`
)

func RetrieveForNotificationsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		select_stmt, err := db.Prepare(SELECT_NOTIFIABLE_STUDENTS)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
			return 
		}

		var notificationRequest models.NotificationRequest

		if err := c.BindJSON(&notificationRequest); err != nil {
			utils.AbortWithBadRequestError(c, RETRIEVE_NOTIFICATIONS_BAD_REQUEST_ERROR_MESSAGE)
			return
		}

		notification_sender := notificationRequest.Teacher
		notification_text := notificationRequest.NotificationText

		mentionedStudents := getMentionedStudents(notification_text)

		var notifiableStudents []string

		rows, err := select_stmt.Query(notification_sender)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var student models.Student
			if err := rows.Scan(&student.Id); err != nil {
				utils.AbortWithInternalServerError(c, err.Error())
				return 
			}
			notifiableStudents = append(notifiableStudents, student.Id)
		}

		notifiableStudents = append(notifiableStudents, mentionedStudents...)

		insert_stmt, err := db.Prepare(INSERT_INTO_NOTIFICATIONS)
		if err != nil {
			utils.AbortWithInternalServerError(c, err.Error())
			return
		}

		for _, student := range notifiableStudents {
			_, err := insert_stmt.Exec(notification_sender, student, notification_text)
			if err != nil {
				utils.AbortWithInternalServerError(c, err.Error())
				return 
			}
		}

		result := map[string][]string{
			"recipients": notifiableStudents,
		}

		c.IndentedJSON(http.StatusOK, result)
	}
}

func getMentionedStudents(text string) []string {
	var mentionedStudents []string
	for _, s := range strings.Split(text, " @")[1:] {
		student := strings.Split(s, " ")[0]
		mentionedStudents = append(mentionedStudents, student)
	}
	return mentionedStudents
}
