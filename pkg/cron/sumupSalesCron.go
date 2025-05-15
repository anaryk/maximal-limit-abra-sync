package cron

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/abra"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/db"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/sumup"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/utils"
)

func PerformSumUpSalesImport(internalDB *db.Connector, abraClient *abra.Connector, sumupClient *sumup.Connector, sumUpMerchant string) {
	// Date exatrction
	now := time.Now().UTC().Truncate(24 * time.Hour)
	now = now.Add(-24 * time.Hour)
	formatted := now.Format("2006-01-02T15:04:05Z")
	// Fetch transactions from SumUp
	transactions, err := sumupClient.GetTransactionInSpecificDate(formatted, sumUpMerchant)
	if err != nil {
		log.Error().Msgf("Error fetching transactions from SumUp: %v", err)
		return
	}

	for _, transaction := range transactions.Items {
		imported, err := internalDB.IsTransactionImported(transaction.ID)
		if err != nil {
			log.Error().Msgf("Error checking transaction import status: %v", err)
			continue
		}

		if transaction.Status != "SUCCESSFUL" {
			log.Debug().Msgf("Transaction %s is not successful, skipping.", transaction.ID)
			continue
		}

		receipt, err := sumupClient.GetReceiptByTransactionID(transaction.ID, sumUpMerchant)
		if err != nil {
			log.Error().Msg(err.Error())
		}

		var items []abra.SaleReceiptItem

		for _, item := range receipt.TransactionData.Products {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				log.Error().Msgf("Error parsing price for item %s: %v", item.Name, err)
				continue
			}
			log.Debug().Msgf("Item: %s, Price: %f", item.Name, price)
			if item.Description != "" {
				eam := utils.ExtractEANCode(item.Description)
				if eam != "" {
					log.Debug().Msgf("EAM code found in description: %s try to import with price tables", eam)
					items = append(items, abra.SaleReceiptItem{
						Cenik:  fmt.Sprintf("code:%s", eam),
						MnozMj: float64(item.Quantity),
						Kod:    eam,
						Sklad:  "code:SKLAD",
					})
					continue
				} else {
					log.Debug().Msgf("No EAM code found in description: %s", item.Description)
				}
			}
			log.Debug().Msgf("No description found for item: %s using non registred field", item.Name)
			items = append(items, abra.SaleReceiptItem{
				Nazev:  item.Name,
				MnozMj: float64(item.Quantity),
				CenaMj: func() float64 {
					totalWithVat, err := strconv.ParseFloat(item.TotalWithVat, 64)
					if err != nil {
						log.Error().Msgf("Error parsing TotalWithVat for item %s: %v", item.Name, err)
						return 0
					}
					return totalWithVat
				}(),
			})

		}

		if imported {
			log.Debug().Msgf("Transaction %s already imported, skipping.", transaction.ID)
			continue
		}

		var prodejka abra.SaleProdejka

		// Set payment method to card due to sumup only accepting card payments
		prodejka.FormaUhradyCis = "code:KARTA"
		//to nwm co znamená
		prodejka.PrimUcet = "code:378001"
		//to nwm co znamená
		prodejka.ProtiUcet = "code:604001"
		//typ ucetní operace
		prodejka.TypUcOp = "code:TRŽBA ZBOŽÍ"
		//typ dokladu
		prodejka.TypDokl = "code:KASA"
		//currency
		prodejka.Mena = "code:CZK"

		// Set the items to the prodejka
		if len(items) == 0 {
			log.Error().Msgf("No items found for transaction %s, skipping.", transaction.ID)
			continue
		}
		// Set the date of the transaction
		prodejka.PolozkyDokladu = items

		salesReceipt := abra.SaleReceipt{
			Winstrom: struct {
				Version  string            `json:"@version"`
				Prodejka abra.SaleProdejka `json:"prodejka"`
			}{
				Version:  "1.0",
				Prodejka: prodejka,
			},
		}

		response, err := abraClient.CreateSalesReceipt(salesReceipt)
		if err != nil {
			log.Printf("Error creating sales receipt in Abra: %v", err)
			continue
		}

		if response.Winstrom.Success == "true" {
			err = internalDB.InsertOrUpdateTransactionState(transaction.ID, true)
			if err != nil {
				log.Error().Msgf("Error updating transaction state in DB: %v", err)
				continue
			}
			log.Info().Msgf("Transaction %s from SumUp imported successfully.", transaction.ID)
		} else {
			log.Error().Msgf("Failed to create sales receipt in Abra: %v", response.Winstrom.Results)
		}
	}
}
