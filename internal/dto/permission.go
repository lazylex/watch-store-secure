package dto

type PermissionDTO struct {
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Description string `json:"description"`
	Service     int    `json:"service"`
}