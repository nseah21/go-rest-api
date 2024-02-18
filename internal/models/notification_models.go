package models

type NotificationRequest struct {
	Teacher string `json:"teacher" binding:"required"`
	NotificationText string `json:"notification" binding:"required"`
}
