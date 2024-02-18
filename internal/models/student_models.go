package models

type Student struct {
	Id string `json:"id"`
}

type SuspensionRequest struct {
	Student string `json:"student"`
}
