package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMessage(message string) map[string]string {
	return map[string]string{
		"message": message,
	}
}

func AbortWithInternalServerError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(
		http.StatusInternalServerError,
		map[string]string {
			"message": message,
		},
	)
}

func AbortWithBadRequestError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(
		http.StatusBadRequest,
		map[string]string {
			"message": message,
		},
	)
}
