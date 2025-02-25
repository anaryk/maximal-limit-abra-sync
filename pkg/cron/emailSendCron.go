package cron

import (
	"github.com/Pacerino/postal-go"
	"github.com/rs/zerolog/log"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/abra"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/db"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/email"
)

func PerformEmailSendCron(internalDB *db.Connector, abraClient *abra.Connector, postalClient *postal.Client) {
	orderToSend, err := internalDB.QueryUnsendOrders()
	if err != nil {
		log.Err(err).Msg("Failed to query unsend orders")
		return
	}
	for _, order := range orderToSend {
		err := email.SendInvoiceEmail(order.Email, order.InvoiceID, postalClient, abraClient)
		if err != nil {
			log.Err(err).Msgf("Failed to send email to %s", order.Email)
			return
		}
		err = internalDB.UpdateEmailSentState(order.OrderNumber)
		if err != nil {
			log.Err(err).Msg("Failed to update email state")
			return
		}
	}

}
