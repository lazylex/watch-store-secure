package dto

type NameServiceSecret struct {
	Name    string `json:"name"`
	Service string `json:"service"`
	Secret  string `json:"secret"`
}
