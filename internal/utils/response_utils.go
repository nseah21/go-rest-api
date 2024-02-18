package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	INTERNAL_SERVER_ERROR_MESSAGE = `Encountered the following error while processing your request... `
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
			"message": INTERNAL_SERVER_ERROR_MESSAGE + message,
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
