package services

import (
	"fmt"
	"log"
	"price-aggregator/models"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Scrapea precios de Walmart
func SearchWalmart(query string) []models.Product {
	var products []models.Product
	const maxProducts = 5
	// Iniciar Colly
	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*walmart.*",
		Parallelism: 2,
		Delay:       500 * time.Millisecond, // Reducimos el delay para mayor velocidad
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		log.Println("Enviando petición a:", r.URL)
	})

	// Extraer información de cada producto
	c.OnHTML("[data-testid='list-view']", func(e *colly.HTMLElement) {
		if len(products) >= maxProducts {
			log.Println("Se ha alcanzado el límite de productos, ignorando más resultados walmart.")
			return // Ignoramos más elementos
		}

		name := e.ChildText("[data-automation-id='product-title']")
		priceWhole := e.ChildText("[data-automation-id='product-price'] .f2")
		priceFraction := e.ChildText("[data-automation-id='product-price'] .f6.f5-l")
		link := e.ChildAttr("a", "href")

		priceFraction = strings.TrimPrefix(priceFraction, "$")

		if name != "" && priceWhole != "" {
			price := fmt.Sprintf("$%s,%s", priceWhole, priceFraction)

			products = append(products, models.Product{
				Name:  name,
				Price: price,
				Link:  "https://www.walmart.com" + link,
				Store: "Walmart",
			})

			if len(products) >= maxProducts {
				c.OnHTMLDetach("[data-testid='list-view']") // Detiene la recolección de más elementos
			}
		}
	})

	// Visitar la URL de búsqueda
	searchURL := fmt.Sprintf("https://www.walmart.com/search/?query=%s", query)
	log.Println("Visitando URL:", searchURL)
	err := c.Visit(searchURL)
	if err != nil {
		log.Println("Error al hacer scraping en Walmart:", err)
	}
	c.OnResponse(func(r *colly.Response) {
		log.Println("Código de respuesta:", r.StatusCode)
		if r.StatusCode != 200 {
			log.Println("La página puede estar bloqueando el scraper o no existe.")
		}
		log.Println("Contenido HTML recibido:")
		log.Println(string(r.Body)) // Esto imprimirá el HTML recibido
	})

	return products
}
