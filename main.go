package main

import (
	"os"

	"github.com/Pacerino/postal-go"
	croner "github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/abra"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/cron"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/db"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/sumup"
)

func main() {
	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// ========================================================================================== //

	if os.Getenv("DB_MAXADMIN_HOST") == "" || os.Getenv("DB_MAXADMIN_USER") == "" || os.Getenv("DB_MAXADMIN_PASSWORD") == "" || os.Getenv("DB_MAXADMIN_NAME") == "" || os.Getenv("DB_INTERNAL_HOST") == "" || os.Getenv("DB_INTERNAL_USER") == "" || os.Getenv("DB_INTERNAL_PASSWORD") == "" || os.Getenv("DB_INTERNAL_NAME") == "" || os.Getenv("ABRA_USER") == "" || os.Getenv("ABRA_PASSWORD") == "" || os.Getenv("POSTAL_URL") == "" || os.Getenv("POSTAL_API_KEY") == "" || os.Getenv("SUMUP_API_TOKEN") == "" || os.Getenv("SUMUP_MERCHANT_ID") == "" {
		log.Fatal().Msg("Missing environment variables")
	}

	if os.Getenv("ENABLE_EMAIL_CRON") == "" {
		os.Setenv("ENABLE_EMAIL_CRON", "false")
	}

	client := postal.NewClient(os.Getenv("POSTAL_URL"), os.Getenv("POSTAL_API_KEY"))

	maxadminDB, err := db.NewMySQLConnector(os.Getenv("DB_MAXADMIN_NAME"), os.Getenv("DB_MAXADMIN_HOST"), os.Getenv("DB_MAXADMIN_USER"), os.Getenv("DB_MAXADMIN_PASSWORD"))
	if err != nil {
		log.Error().Msg(err.Error())
	}

	intertnalDB, err := db.NewMySQLConnector(os.Getenv("DB_INTERNAL_NAME"), os.Getenv("DB_INTERNAL_HOST"), os.Getenv("DB_INTERNAL_USER"), os.Getenv("DB_INTERNAL_PASSWORD"))
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// Check if internal tables are populated and create them if not
	intertnalDB.InitInternalDBIfNotExist()
	intertnalDB.InitSumupTransactionStateTable()
	// END table creation

	abraClient := abra.NewAbraConnector(os.Getenv("ABRA_USER"), os.Getenv("ABRA_PASSWORD"))

	sumClient := sumup.NewSumUpAPI(os.Getenv("SUMUP_API_TOKEN"))

	//Run crons on start container
	cron.PerformOrderInvoiceSync(maxadminDB, intertnalDB, abraClient)
	cron.PerformTicketsInvoiceSync(maxadminDB, intertnalDB, abraClient)
	cron.PerformChipInvoiceSync(maxadminDB, intertnalDB, abraClient)
	if os.Getenv("ENABLE_EMAIL_CRON") == "true" {
		cron.PerformEmailSendCron(intertnalDB, abraClient, client)
	}
	cron.PerformSumUpSalesImport(intertnalDB, abraClient, sumClient, os.Getenv("SUMUP_MERCHANT_ID"))
	//init cron
	c := croner.New()

	// Add PerformOrderInvoiceSync job to run every 4 hours
	_, err = c.AddFunc("@every 4h", func() {
		cron.PerformOrderInvoiceSync(maxadminDB, intertnalDB, abraClient)
		cron.PerformTicketsInvoiceSync(maxadminDB, intertnalDB, abraClient)
		cron.PerformChipInvoiceSync(maxadminDB, intertnalDB, abraClient)
		if os.Getenv("ENABLE_EMAIL_CRON") == "true" {
			cron.PerformEmailSendCron(intertnalDB, abraClient, client)
		}
		cron.PerformSumUpSalesImport(intertnalDB, abraClient, sumClient, os.Getenv("SUMUP_MERCHANT_ID"))
	})
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// Add PerformXMLFeedSync job to run every 7 days
	_, err = c.AddFunc("@every 168h", func() {
		cron.PerformXMLFeedSync(abraClient)
	})
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// Start the cron scheduler
	c.Start()

	// Prevent the program from exiting
	select {}
}
