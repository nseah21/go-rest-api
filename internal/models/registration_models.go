package models

type RegistrationRequest struct {
	Teacher string `json:"teacher" binding:"required"`
	Students []string `json:"students" binding:"required"`
}
