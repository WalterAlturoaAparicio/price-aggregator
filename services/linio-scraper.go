package services

import (
	"fmt"
	"log"
	"price-aggregator/models"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// SearchLinio hace scraping en Linio
func SearchLinio(query string) []models.Product {
	var products []models.Product
	const maxProducts = 5

	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*linio.*",
		Parallelism: 2,
		Delay:       500 * time.Millisecond, // Reducimos el delay para mayor velocidad
	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Referer", "https://www.google.com/")
		log.Println("Enviando petición a:", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		log.Println("Código de respuesta:", r.StatusCode)
		if r.StatusCode != 200 {
			log.Println("La página puede estar bloqueando el scraper o no existe.")
		}
	})

	c.OnHTML(".grid-pod", func(e *colly.HTMLElement) {
		if len(products) >= maxProducts {
			log.Println("Se ha alcanzado el límite de productos, ignorando más resultados. linio")
			return // Ignoramos más elementos
		}
		name := e.ChildText(".pod-subTitle")
		price := e.ChildText(".prices-1 span") // Precio normal
		link := e.Request.AbsoluteURL(e.ChildAttr(".pod-link", "href"))

		if name != "" && price != "" {
			products = append(products, models.Product{
				Name:  name,
				Price: price,
				Link:  link,
				Store: "Linio",
			})
			if len(products) >= maxProducts {
				c.OnHTMLDetach(".grid-pod") // Detiene la recolección de más elementos
			}
		}
	})

	// Visitar la URL de búsqueda
	searchURL := fmt.Sprintf("https://linio.falabella.com.co/linio-co/search?Ntt=%s", strings.ReplaceAll(query, " ", "+"))
	log.Println("Visitando URL:", searchURL)
	err := c.Visit(searchURL)
	if err != nil {
		log.Println("Error al hacer scraping en Linio:", err)
	}

	return products
}
