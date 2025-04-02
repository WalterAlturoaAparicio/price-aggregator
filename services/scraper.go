package services

import (
	"price-aggregator/models"
	"sync"
)

// SearchProducts consulta múltiples fuentes en paralelo
func SearchProducts(query string) []models.Product {
	var wg sync.WaitGroup
	results := make(chan models.Product, 30)

	sources := []func(string) []models.Product{
		SearchLinio,
		SearchWalmart,
		SearchMercadoLibre,
	}

	// Ejecutar cada búsqueda en paralelo
	for _, source := range sources {
		wg.Add(1)
		go func(src func(string) []models.Product) {
			defer wg.Done()
			for _, product := range src(query) {
				results <- product
			}
		}(source)
	}

	// Esperar a que todas las goroutines terminen
	go func() {
		wg.Wait()
		close(results)
	}()

	// Convertir canal a slice
	var products []models.Product
	for p := range results {
		products = append(products, p)
	}

	return products
}
