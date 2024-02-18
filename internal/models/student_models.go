package models

type Student struct {
	Id string `json:"id" binding:"required"`
}

type SuspensionRequest struct {
	Student string `json:"student" binding:"required"`
}
