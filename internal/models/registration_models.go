package models

type RegistrationRequest struct {
	Teacher string `json:"teacher"`
	Students []string `json:"students"`
}
