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
	addedCount := 0
	fixedCount := 0
	correctCount := 0
	skippedCount := 0

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

		// Check if this order has a voucher in the database
		voucher, err := maxadminDB.QueryVoucherByOrderID(ticket.OrderID)
		if err != nil {
			log.Err(err).Msgf("Failed to query voucher for order %s", order.OrderNumber)
			continue
		}
		if voucher == nil {
			skippedCount++
			continue
		}

		// Calculate the correct voucher amount
		correctAmount := utils.CalculateTotalPriceWithVat(voucher.Price, float64(voucher.Vat))

		// Get existing voucher items on the invoice with their amounts
		existingVoucherItems, err := abraClient.GetInvoiceVoucherItemsWithAmount(ticket.InvoiceNum)
		if err != nil {
			log.Err(err).Msgf("Failed to check voucher items for invoice %s", ticket.InvoiceNum)
			continue
		}

		// Case 1: No voucher item - need to add
		if len(existingVoucherItems) == 0 {
			item := abra.FakturaPolozka{
				Popis:   fmt.Sprintf("Slevový poukaz %s", voucher.VoucherCode),
				Pocet:   1,
				CenaKus: correctAmount,
			}
			resp, err := abraClient.AddInvoiceItem(ticket.InvoiceNum, item)
			if err != nil {
				log.Err(err).Msgf("Failed to add voucher item to invoice %s", ticket.InvoiceNum)
				continue
			}
			if resp.Winstrom.Success == "true" {
				addedCount++
				log.Info().Msgf("Added missing voucher to invoice %s (discount: %.2f)", ticket.InvoiceNum, correctAmount)
			}
			continue
		}

		// Case 2: Exactly one voucher item - check if amount is correct
		if len(existingVoucherItems) == 1 {
			currentAmount := existingVoucherItems[0].Amount
			if currentAmount == correctAmount {
				correctCount++
				log.Debug().Msgf("Invoice %s has correct voucher amount %.2f", ticket.InvoiceNum, correctAmount)
				continue
			}
			// Amount is wrong - remove and re-add
			log.Info().Msgf("Invoice %s has wrong voucher amount %.2f, should be %.2f - fixing", ticket.InvoiceNum, currentAmount, correctAmount)
		}

		// Case 3: Multiple voucher items OR wrong amount - remove all and add correct one
		for _, existingItem := range existingVoucherItems {
			if err := abraClient.RemoveInvoiceItem(ticket.InvoiceNum, existingItem.ID); err != nil {
				log.Err(err).Msgf("Failed to remove voucher item %s from invoice %s", existingItem.ID, ticket.InvoiceNum)
			}
		}

		// Add correct voucher item
		item := abra.FakturaPolozka{
			Popis:   fmt.Sprintf("Slevový poukaz %s", voucher.VoucherCode),
			Pocet:   1,
			CenaKus: correctAmount,
		}
		resp, err := abraClient.AddInvoiceItem(ticket.InvoiceNum, item)
		if err != nil {
			log.Err(err).Msgf("Failed to add voucher item to invoice %s", ticket.InvoiceNum)
			continue
		}
		if resp.Winstrom.Success == "true" {
			fixedCount++
			log.Info().Msgf("Fixed invoice %s voucher amount to %.2f", ticket.InvoiceNum, correctAmount)
		} else {
			log.Warn().Msgf("Failed to fix invoice %s: %v", ticket.InvoiceNum, resp.Winstrom.Results)
		}
	}

	log.Info().Msgf("Invoice repair completed: %d added, %d fixed (wrong amount), %d already correct, %d skipped (no voucher)", addedCount, fixedCount, correctCount, skippedCount)
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
			{Popis: fmt.Sprintf("Fakturujeme vám nákup kreditů (%.0f) dle objednávky %s ze dne %s", credit.Count, credit.OrderNumber, credit.Created), Pocet: 1, CenaKus: utils.CalculateTotalPriceWithVat(credit.TotalPrice, float64(credit.Vat))},
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
