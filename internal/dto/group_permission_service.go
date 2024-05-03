package dto

type GroupPermissionService struct {
	Group      string `json:"group"`
	Permission string `json:"permission"`
	Service    string `json:"service"`
}
