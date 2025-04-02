package models

// Product representa un resultado de búsqueda
type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Link  string `json:"link"`
	Store string `json:"store"`
}
