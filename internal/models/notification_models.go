package models

type NotificationRequest struct {
	Teacher string `json:"teacher"`
	NotificationText string `json:"notification"`
}
