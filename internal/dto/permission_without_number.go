package dto

type PermissionWithoutNumberDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Service     string `json:"service"`
}
