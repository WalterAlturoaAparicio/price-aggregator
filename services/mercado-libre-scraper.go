package services

import (
	"fmt"
	"log"
	"price-aggregator/models"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// SearchMercadoLibre hace scraping en Mercado Libre
func SearchMercadoLibre(query string) []models.Product {
	var products []models.Product
	const maxProducts = 5
	// Iniciar Colly
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		log.Println("Enviando petición a:", r.URL)
	})
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*mercadolibre.*",
		Parallelism: 2,
		Delay:       500 * time.Millisecond, // Reducimos el delay para mayor velocidad
	})

	// Extraer información de cada producto
	c.OnHTML(".ui-search-layout__item", func(e *colly.HTMLElement) {
		if len(products) >= maxProducts {
			log.Println("Se ha alcanzado el límite de productos, ignorando más resultados. ML")
			return // Ignoramos más elementos
		}

		name := e.ChildText(".poly-component__title")
		price := e.ChildText(".poly-price__current .andes-money-amount__fraction")
		currency := e.ChildText(".poly-price__current .andes-money-amount__currency-symbol")
		link := e.ChildAttr(".poly-component__title-wrapper a", "href")

		if name != "" && price != "" {
			priceFormatted := fmt.Sprintf("%s%s", currency, price)
			if !strings.HasPrefix(link, "http") {
				link = "https://www.mercadolibre.com.co" + link
			}

			products = append(products, models.Product{
				Name:  name,
				Price: priceFormatted,
				Link:  link,
				Store: "Mercado Libre",
			})

			if len(products) >= maxProducts {
				c.OnHTMLDetach(".ui-search-layout__item") // Detiene la recolección de más elementos
			}
		}
	})

	// Construir la URL de búsqueda
	searchURL := fmt.Sprintf("https://www.mercadolibre.com.co/jm/search?as_word=%s", strings.ReplaceAll(query, " ", "+"))
	log.Println("Visitando URL:", searchURL)
	err := c.Visit(searchURL)
	if err != nil {
		log.Println("Error al hacer scraping en Mercado Libre:", err)
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
