package email

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/Pacerino/postal-go"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/abra"
)

func SendInvoiceEmail(email string, invoiceID string, postalClien *postal.Client, abraClient *abra.Connector) error {
	invBase64, err := abraClient.GetPDFInvoiceAsBase64(invoiceID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to download invoice")
	}
	var attachments []Attachment
	attachments = append(attachments, Attachment{Name: fmt.Sprintf("faktura-%s.pdf", invoiceID), Content: "application/pdf", Base64Data: invBase64})
	message := &postal.SendRequest{
		To:      []string{email},
		From:    "info@maximal-limit.cz",
		Subject: fmt.Sprintf("Faktura %s za služby Maximal Limit - Již uhrazeno ! Nehraďte - ", invoiceID),
		HTMLBody: `<html>
            <body>
                <p>Dobrý den,</p>
                <p>Pořádek dělá přátele a proto zasíláme fakturu za naše služby.</p>
                <p>Pokud máte jakékoli dotazy, neváhejte nás kontaktovat.</p>
				<p><strong>Faktura byla hrazena platební kartou a je již uhrazena.</strong></p>
                <p>S pozdravem,</p>
                <p><strong>Maximal Limit</strong></p>
            </body>
        </html>`,
		Attachments: attachments,
	}
	resp, _, err := postalClien.Send.Send(context.TODO(), message)
	if err != nil {
		log.Err(err).Msg("Failed to send email")
	}
	log.Info().Msgf("Email sent with message ID: %s", resp.MessageID)
	return nil
}
