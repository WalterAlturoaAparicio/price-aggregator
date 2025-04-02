package handlers

import (
	"price-aggregator/services"

	"github.com/gofiber/fiber/v2"
)

// SearchHandler maneja las solicitudes de búsqueda
func SearchHandler(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Debes proporcionar un término de búsqueda"})
	}

	products := services.SearchProducts(query)

	return c.JSON(products)
}
