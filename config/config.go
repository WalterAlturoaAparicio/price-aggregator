package config

import (
	"os"
)

// eBay API Key (debes configurarla como variable de entorno)
var EbayAPIKey = os.Getenv("EBAY_API_KEY")
