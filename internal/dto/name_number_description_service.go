package dto

type NameNumberDescriptionService struct {
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Description string `json:"description"`
	Service     string `json:"service"`
}
