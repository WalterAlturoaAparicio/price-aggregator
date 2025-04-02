package services

import (
	"fmt"
	"log"
	"math/rand"
	"price-aggregator/models"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// SearchAmazon hace scraping en Amazon
func SearchAmazon(query string) []models.Product {
	var products []models.Product
	const maxProducts = 5
	var userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
	}
	// Iniciar Colly
	c := colly.NewCollector(
		colly.AllowedDomains("www.amazon.com", "amazon.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
		r.Headers.Set("Referer", "https://www.google.com/")
		log.Println("Enviando petición a:", r.URL)
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*amazon.*",
		Parallelism: 2,
		Delay:       500 * time.Millisecond, // Reducimos el delay para mayor velocidad
	})

	// Extraer información de cada producto
	c.OnHTML("[data-component-type='s-search-result']", func(e *colly.HTMLElement) {
		if len(products) >= maxProducts {
			log.Println("Se ha alcanzado el límite de productos, ignorando más resultados. Amazon")
			return // Ignoramos más elementos
		}

		name := e.ChildText("[data-cy='title-recipe'] span")
		price := e.ChildText("[data-cy='secondary-offer-recipe'] .a-color-base")
		link := e.ChildAttr(".aok-relative a", "href")

		if name != "" && price != "" {
			priceFormatted := fmt.Sprintf("%s", price)
			if !strings.HasPrefix(link, "http") {
				link = "https://www.amazon.com" + link
			}

			products = append(products, models.Product{
				Name:  name,
				Price: priceFormatted,
				Link:  link,
				Store: "Amazon",
			})

			if len(products) >= maxProducts {
				c.OnHTMLDetach("[data-component-type='s-search-result']") // Detiene la recolección de más elementos
			}
		}
	})

	// Construir la URL de búsqueda
	searchURL := fmt.Sprintf("https://www.amazon.com/s?k=%s", strings.ReplaceAll(query, " ", "+"))
	log.Println("Visitando URL:", searchURL)
	err := c.Visit(searchURL)
	if err != nil {
		log.Println("Error al hacer scraping en Amazon:", err)
	}

	c.OnResponse(func(r *colly.Response) {
		log.Println("Código de respuesta:", r.StatusCode)
		if r.StatusCode == 503 {
			log.Println("Amazon está bloqueando la solicitud. Intenta con un proxy o espera.")
			return
		}
	})

	return products
}
