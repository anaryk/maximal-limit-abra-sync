package cron

import (
	"html"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/abra"
	xmlfeeds "github.com/anaryk/maximal-limit-abra-sync/pkg/xml_feeds"
	"github.com/rs/zerolog/log"
)

func PerformXMLFeedSync(abraClient *abra.Connector) {
	nutrend, err := xmlfeeds.GetNutrendFeed()
	if err != nil {
		log.Error().Msgf("Error fetching XML feed: %v", err)
		return
	}

	dafit, err := xmlfeeds.GetDafitFeed()
	if err != nil {
		log.Error().Msgf("Error fetching XML feed: %v", err)
		return
	}

	for _, product := range nutrend.ShopItem {
		log.Debug().Msgf("Product: %s, Price: %s, EAM: %s, Description: %s", product.ProductName, product.PriceVAT, product.EAN, product.Description)
		abraProduct := abra.Cenik{
			EanKod:        product.EAN,
			Kod:           product.EAN,
			Nazev:         product.ProductName,
			NakupCena:     product.PriceVAT,
			Popis:         html.UnescapeString(product.Description),
			ExportNaEshop: "false",
			Dodavatel:     "code:NUTREND",
			Skladove:      "true",
			CenJednotka:   "1.0",
			EvidExpir:     "true",
			Mj1:           "code:KS",
			ProdejKasa:    "true",
			SkupZboz:      "code:ZBOŽÍ",
		}
		exists, err := abraClient.CheckIfPriceItemExists(product.EAN)
		if err != nil {
			log.Error().Msgf("Error checking if price item exists: %v", err)
			continue
		}
		if exists {
			log.Info().Msgf("Price item %s already exists, skipping creation.", product.EAN)
			continue
		}
		log.Info().Msgf("Creating price item %s", product.EAN)
		response, err := abraClient.CreatePriceItem(abraProduct)
		if err != nil {
			log.Error().Msgf("Error creating price item: %v", err)
			continue
		}
		if response.Winstrom.Success != "true" {
			log.Error().Msgf("Error creating price item: %s", response.Winstrom.Results)
			continue
		}
		log.Info().Msgf("Price item %s created successfully", product.EAN)
	}

	for _, product := range dafit.ShopItem {
		log.Debug().Msgf("Product: %s, Price: %s, EAM: %s, Description: %s", product.ProductName, product.PriceVAT, product.EAN, product.Description)
		abraProduct := abra.Cenik{
			EanKod:        product.EAN,
			Kod:           product.EAN,
			Nazev:         product.ProductName,
			NakupCena:     product.PriceVAT,
			Popis:         html.UnescapeString(product.Description),
			ExportNaEshop: "false",
			Dodavatel:     "code:DAFIT",
			Skladove:      "true",
			CenJednotka:   "1.0",
			EvidExpir:     "true",
			Mj1:           "code:KS",
			ProdejKasa:    "true",
			SkupZboz:      "code:ZBOŽÍ",
		}
		exists, err := abraClient.CheckIfPriceItemExists(product.EAN)
		if err != nil {
			log.Error().Msgf("Error checking if price item exists: %v", err)
			continue
		}
		if exists {
			log.Info().Msgf("Price item %s already exists, skipping creation.", product.EAN)
			continue
		}
		log.Info().Msgf("Creating price item %s", product.EAN)
		response, err := abraClient.CreatePriceItem(abraProduct)
		if err != nil {
			log.Error().Msgf("Error creating price item: %v", err)
			continue
		}
		if response.Winstrom.Success != "true" {
			log.Error().Msgf("Error creating price item: %s", response.Winstrom.Results)
			continue
		}
		log.Info().Msgf("Price item %s created successfully", product.EAN)
	}

}
