package cron

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/abra"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/db"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/utils"
)

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
		resp, err := abraClient.CreateInvoice(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)), utils.GetCurrentDate(), utils.GetCurrentDate(), order.InvoiceNum, items)
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
		resp, err := abraClient.CreateInvoice(utils.GenerateShortCode(fmt.Sprintf("%s %s", user.Name, user.Surname)), utils.GetCurrentDate(), utils.GetCurrentDate(), ticket.InvoiceNum, items)
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
