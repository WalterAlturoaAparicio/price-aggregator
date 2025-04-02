package main

import (
	"log"

	"price-aggregator/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Ruta de búsqueda
	app.Get("/search", handlers.SearchHandler)

	log.Println("Servidor corriendo en http://localhost:3000")
	app.Listen(":3000")
}
