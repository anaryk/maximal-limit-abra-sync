package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/abra"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/db"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/utils"
)

func RepairTicketInvoicesWithMissingVouchers(maxadminDB, internalDB *db.Connector, abraClient *abra.Connector) {
	importedOrders, err := internalDB.QueryAllImportedOrders()
	if err != nil {
		log.Err(err).Msg("Failed to query imported orders")
		return
	}

	currentYear := fmt.Sprintf("%d", time.Now().Year())
	repairedCount := 0
	duplicatesRemovedCount := 0
	skippedCount := 0
	alreadyHasVoucherCount := 0

	for _, order := range importedOrders {
		ticket, err := maxadminDB.QueryTicketByOrderNumber(order.OrderNumber)
		if err != nil {
			log.Err(err).Msgf("Failed to query ticket for order %s", order.OrderNumber)
			continue
		}
		if ticket == nil {
			skippedCount++
			continue
		}

		// Filter only current year invoices (invoice codes like P2025xxxx or P2026xxxx)
		if !strings.Contains(ticket.InvoiceNum, currentYear) {
			skippedCount++
			continue
		}

		// First, remove any duplicate voucher items
		removed, err := abraClient.RemoveDuplicateVoucherItems(ticket.InvoiceNum)
		if err != nil {
			log.Err(err).Msgf("Failed to remove duplicates from invoice %s", ticket.InvoiceNum)
		}
		if removed > 0 {
			duplicatesRemovedCount += removed
			log.Info().Msgf("Removed %d duplicate voucher items from invoice %s", removed, ticket.InvoiceNum)
		}

		// Check if invoice already has a voucher item
		hasVoucher, err := abraClient.HasVoucherItem(ticket.InvoiceNum)
		if err != nil {
			log.Err(err).Msgf("Failed to check voucher items for invoice %s", ticket.InvoiceNum)
			continue
		}
		if hasVoucher {
			alreadyHasVoucherCount++
			log.Debug().Msgf("Invoice %s already has voucher item, skipping", ticket.InvoiceNum)
			continue
		}

		voucher, err := maxadminDB.QueryVoucherByOrderID(ticket.OrderID)
		if err != nil {
			log.Err(err).Msgf("Failed to query voucher for order %s", order.OrderNumber)
			continue
		}
		if voucher == nil {
			skippedCount++
			continue
		}

		item := abra.FakturaPolozka{
			Popis:   fmt.Sprintf("Slevový poukaz %s", voucher.VoucherCode),
			Pocet:   1,
			CenaKus: utils.CalculateTotalPriceWithVat(voucher.Price, float64(voucher.Vat)),
		}

		resp, err := abraClient.AddInvoiceItem(ticket.InvoiceNum, item)
		if err != nil {
			log.Err(err).Msgf("Failed to add voucher item to invoice %s", ticket.InvoiceNum)
			continue
		}

		if resp.Winstrom.Success == "true" {
			repairedCount++
			log.Info().Msgf("Repaired invoice %s with voucher %s (discount: %.2f)", ticket.InvoiceNum, voucher.VoucherCode, voucher.Price)
		} else {
			log.Warn().Msgf("Failed to repair invoice %s: %v", ticket.InvoiceNum, resp.Winstrom.Results)
		}
	}

	log.Info().Msgf("Invoice repair completed: %d repaired, %d duplicates removed, %d already had voucher, %d skipped", repairedCount, duplicatesRemovedCount, alreadyHasVoucherCount, skippedCount)
}

func PerformOrderInvoiceSync(maxadminDB, internalDB *db.Connector, abraClient *abra.Connector) {
	orderData, err := maxadminDB.QueryPayedOrdersInYear(utils.GetFirstDayOfActualYear())
	if err != nil {
		log.Err(err).Msg("Failed to query payed orders")
		return
	}
	for _, order := range orderData {
		state, err := internalDB.QueryOrderProccesedState(order.OrderNumber)
		if err != nil {
			log.Err(err).Msg("Failed to query order state")
			return
		}
		if state == internal.InternalDBStatusImported || state != "" {
			log.Debug().Msg("Order already imported")
			continue
		}
		user, err := maxadminDB.QueryUserByID(order.UserID)
		if err != nil {
			log.Err(err).Msg("Failed to query user")
			return
		}
		contactExist, err := abraClient.CheckIfContactExist(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)))
		if err != nil {
			log.Err(err).Msg("Failed to check if contact exist")
			return
		}
		if len(contactExist.Winstrom.Adresar) == 0 {
			contact := abra.ContactData{
				Name:       fmt.Sprintf("%s %s", user.Name, user.Surname),
				Street:     fmt.Sprintf("%s %s", user.Street, user.HouseNumber),
				City:       user.City,
				PostalCode: user.ZipCode,
				Email:      user.Email,
				Mobile:     user.Phone,
			}
			_, err := abraClient.CreateContact(contact)
			if err != nil {
				log.Err(err).Msg("Failed to create contact")
				return
			}
		}
		items := []abra.FakturaPolozka{
			{Popis: fmt.Sprintf("Fakturujeme vám služby dle objednávky %s ze dne %s", order.OrderNumber, order.Created), Pocet: 1, CenaKus: utils.CalculateTotalPriceWithVat(order.TotalPrice, float64(order.Vat))},
		}
		resp, err := abraClient.CreateInvoice(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)), utils.ExtractDate(order.Created), utils.ExtractDate(order.Created), order.InvoiceNum, items)
		if err != nil {
			log.Err(err).Msg("Failed to create invoice")
			return
		}
		if resp.Winstrom.Success == "true" {
			err := internalDB.InsertOrUpdateProcessedState(order.OrderNumber, internal.InternalDBStatusImported, user.Email, resp.Winstrom.Results[0].ID)
			if err != nil {
				log.Err(err).Msg("Failed to insert order status")
				return
			}
			log.Info().Msgf("Order %s imported: %s", order.OrderNumber, resp.Winstrom.Results)
		}
	}
}

func PerformTicketsInvoiceSync(maxadminDB, internalDB *db.Connector, abraClient *abra.Connector) {
	ticketData, err := maxadminDB.QueryPayedTicketsInYear(utils.GetFirstDayOfActualYear())
	if err != nil {
		log.Err(err).Msg("Failed to query payed tickets")
		return
	}
	for _, ticket := range ticketData {
		state, err := internalDB.QueryOrderProccesedState(ticket.OrderNumber)
		if err != nil {
			log.Err(err).Msg("Failed to query ticket state")
			return
		}
		if state == internal.InternalDBStatusImported || state != "" {
			log.Debug().Msg("Ticket already imported")
			continue
		}
		user, err := maxadminDB.QueryUserByID(ticket.UserID)
		if err != nil {
			log.Err(err).Msg("Failed to query user")
			return
		}
		contactExist, err := abraClient.CheckIfContactExist(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)))
		if err != nil {
			log.Err(err).Msg("Failed to check if contact exist")
			return
		}
		if len(contactExist.Winstrom.Adresar) == 0 {
			contact := abra.ContactData{
				Name:       fmt.Sprintf("%s %s", user.Name, user.Surname),
				Street:     fmt.Sprintf("%s %s", user.Street, user.HouseNumber),
				City:       user.City,
				PostalCode: user.ZipCode,
				Email:      user.Email,
				Mobile:     user.Phone,
			}
			_, err := abraClient.CreateContact(contact)
			if err != nil {
				log.Err(err).Msg("Failed to create contact")
				return
			}
		}
		items := []abra.FakturaPolozka{
			{Popis: fmt.Sprintf("Fakturujeme vám permanentku %s ze dne %s", ticket.OrderNumber, ticket.Created), Pocet: 1, CenaKus: utils.CalculateTotalPriceWithVat(ticket.TotalPrice, float64(ticket.Vat))},
		}

		voucher, err := maxadminDB.QueryVoucherByOrderID(ticket.OrderID)
		if err != nil {
			log.Err(err).Msg("Failed to query voucher")
			return
		}
		if voucher != nil {
			items = append(items, abra.FakturaPolozka{
				Popis:   fmt.Sprintf("Slevový poukaz %s", voucher.VoucherCode),
				Pocet:   1,
				CenaKus: utils.CalculateTotalPriceWithVat(voucher.Price, float64(voucher.Vat)),
			})
		}

		resp, err := abraClient.CreateInvoice(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)), utils.ExtractDate(ticket.Created), utils.ExtractDate(ticket.Created), ticket.InvoiceNum, items)
		if err != nil {
			log.Err(err).Msg("Failed to create invoice")
			return
		}
		if resp.Winstrom.Success == "true" {
			err := internalDB.InsertOrUpdateProcessedState(ticket.OrderNumber, internal.InternalDBStatusImported, user.Email, resp.Winstrom.Results[0].ID)
			if err != nil {
				log.Err(err).Msg("Failed to insert ticket status")
				return
			}
			log.Info().Msgf("Ticket %s imported: %s", ticket.OrderNumber, resp.Winstrom.Results)
		}
	}

}

func PerformChipInvoiceSync(maxadminDB, internalDB *db.Connector, abraClient *abra.Connector) {
	chipData, err := maxadminDB.QueryPaysChipsInYear(utils.GetFirstDayOfActualYear())
	if err != nil {
		log.Err(err).Msg("Failed to query payed chips")
		return
	}
	for _, chip := range chipData {
		state, err := internalDB.QueryOrderProccesedState(chip.OrderNumber)
		if err != nil {
			log.Err(err).Msg("Failed to query chip state")
			return
		}
		if state == internal.InternalDBStatusImported || state != "" {
			log.Debug().Msg("Chip already imported")
			continue
		}
		user, err := maxadminDB.QueryUserByID(chip.UserID)
		if err != nil {
			log.Err(err).Msg("Failed to query user")
			return
		}
		contactExist, err := abraClient.CheckIfContactExist(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)))
		if err != nil {
			log.Err(err).Msg("Failed to check if contact exist")
			return
		}
		if len(contactExist.Winstrom.Adresar) == 0 {
			contact := abra.ContactData{
				Name:       fmt.Sprintf("%s %s", user.Name, user.Surname),
				Street:     fmt.Sprintf("%s %s", user.Street, user.HouseNumber),
				City:       user.City,
				PostalCode: user.ZipCode,
				Email:      user.Email,
				Mobile:     user.Phone,
			}
			_, err := abraClient.CreateContact(contact)
			if err != nil {
				log.Err(err).Msg("Failed to create contact")
				return
			}
		}
		items := []abra.FakturaPolozka{
			{Popis: fmt.Sprintf("Fakturujeme vám čipy dle objednávky %s ze dne %s", chip.OrderNumber, chip.Created), Pocet: 1, CenaKus: utils.CalculateTotalPriceWithVat(chip.TotalPrice, float64(chip.VAT))},
		}
		resp, err := abraClient.CreateInvoice(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)), utils.ExtractDate(chip.Created), utils.ExtractDate(chip.Created), chip.InvoiceNum, items)
		if err != nil {
			log.Err(err).Msg("Failed to create invoice")
			return
		}
		if resp.Winstrom.Success == "true" {
			err := internalDB.InsertOrUpdateProcessedState(chip.OrderNumber, internal.InternalDBStatusImported, user.Email, resp.Winstrom.Results[0].ID)
			if err != nil {
				log.Err(err).Msg("Failed to insert chip status")
				return
			}
			log.Info().Msgf("Chip %s imported: %s", chip.OrderNumber, resp.Winstrom.Results)
		}
	}
}

func PerformCreditInvoiceSync(maxadminDB, internalDB *db.Connector, abraClient *abra.Connector) {
	creditData, err := maxadminDB.QueryCreditOrdersInYear(utils.GetFirstDayOfActualYear())
	if err != nil {
		log.Err(err).Msg("Failed to query credit orders")
		return
	}
	for _, credit := range creditData {
		state, err := internalDB.QueryOrderProccesedState(credit.OrderNumber)
		if err != nil {
			log.Err(err).Msg("Failed to query credit state")
			return
		}
		if state == internal.InternalDBStatusImported || state != "" {
			log.Debug().Msg("Credit already imported")
			continue
		}
		user, err := maxadminDB.QueryUserByID(credit.UserID)
		if err != nil {
			log.Err(err).Msg("Failed to query user")
			return
		}
		contactExist, err := abraClient.CheckIfContactExist(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)))
		if err != nil {
			log.Err(err).Msg("Failed to check if contact exist")
			return
		}
		if len(contactExist.Winstrom.Adresar) == 0 {
			contact := abra.ContactData{
				Name:       fmt.Sprintf("%s %s", user.Name, user.Surname),
				Street:     fmt.Sprintf("%s %s", user.Street, user.HouseNumber),
				City:       user.City,
				PostalCode: user.ZipCode,
				Email:      user.Email,
				Mobile:     user.Phone,
			}
			_, err := abraClient.CreateContact(contact)
			if err != nil {
				log.Err(err).Msg("Failed to create contact")
				return
			}
		}
		items := []abra.FakturaPolozka{
			{Popis: fmt.Sprintf("Fakturujeme vám nákup kreditů (%s) dle objednávky %s ze dne %s", fmt.Sprintf("%d", credit.Count), credit.OrderNumber, credit.Created), Pocet: 1, CenaKus: utils.CalculateTotalPriceWithVat(credit.TotalPrice, float64(credit.Vat))},
		}
		resp, err := abraClient.CreateInvoice(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)), utils.ExtractDate(credit.Created), utils.ExtractDate(credit.Created), credit.InvoiceNum, items)
		if err != nil {
			log.Err(err).Msg("Failed to create invoice")
			return
		}
		if resp.Winstrom.Success == "true" {
			err := internalDB.InsertOrUpdateProcessedState(credit.OrderNumber, internal.InternalDBStatusImported, user.Email, resp.Winstrom.Results[0].ID)
			if err != nil {
				log.Err(err).Msg("Failed to insert credit status")
				return
			}
			log.Info().Msgf("Credit %s imported: %s", credit.OrderNumber, resp.Winstrom.Results)
		}
	}
}
