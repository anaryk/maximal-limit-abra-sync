package sumup

import "time"

type SumUpTransaction struct {
	Items []struct {
		Amount              float64   `json:"amount,omitempty"`
		CardType            string    `json:"card_type,omitempty"`
		ClientTransactionID string    `json:"client_transaction_id,omitempty"`
		Currency            string    `json:"currency,omitempty"`
		EntryMode           string    `json:"entry_mode,omitempty"`
		ID                  string    `json:"id,omitempty"`
		InstallmentsCount   int       `json:"installments_count,omitempty"`
		PaymentType         string    `json:"payment_type,omitempty"`
		PayoutPlan          string    `json:"payout_plan,omitempty"`
		PayoutsReceived     int       `json:"payouts_received,omitempty"`
		PayoutsTotal        int       `json:"payouts_total,omitempty"`
		ProductSummary      string    `json:"product_summary,omitempty"`
		RefundedAmount      int       `json:"refunded_amount,omitempty"`
		Status              string    `json:"status,omitempty"`
		Timestamp           time.Time `json:"timestamp,omitempty"`
		TransactionCode     string    `json:"transaction_code,omitempty"`
		TransactionID       string    `json:"transaction_id,omitempty"`
		Type                string    `json:"type,omitempty"`
		User                string    `json:"user,omitempty"`
	} `json:"items,omitempty"`
	Links []struct {
		Href string `json:"href,omitempty"`
		Rel  string `json:"rel,omitempty"`
	} `json:"links,omitempty"`
}

type SumUpReceipts struct {
	TransactionData struct {
		TransactionCode    string    `json:"transaction_code,omitempty"`
		TransactionID      string    `json:"transaction_id,omitempty"`
		Amount             string    `json:"amount,omitempty"`
		VatAmount          string    `json:"vat_amount,omitempty"`
		TipAmount          string    `json:"tip_amount,omitempty"`
		Currency           string    `json:"currency,omitempty"`
		Timestamp          time.Time `json:"timestamp,omitempty"`
		Status             string    `json:"status,omitempty"`
		PaymentType        string    `json:"payment_type,omitempty"`
		EntryMode          string    `json:"entry_mode,omitempty"`
		VerificationMethod string    `json:"verification_method,omitempty"`
		CardReader         struct {
			Code string `json:"code,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"card_reader,omitempty"`
		Card struct {
			Last4Digits string `json:"last_4_digits,omitempty"`
			Type        string `json:"type,omitempty"`
			Token       string `json:"token,omitempty"`
		} `json:"card,omitempty"`
		InstallmentsCount int `json:"installments_count,omitempty"`
		Products          []struct {
			Name            string `json:"name,omitempty"`
			Description     string `json:"description,omitempty"`
			Price           string `json:"price,omitempty"`
			VatRate         string `json:"vat_rate,omitempty"`
			SingleVatAmount string `json:"single_vat_amount,omitempty"`
			PriceWithVat    string `json:"price_with_vat,omitempty"`
			VatAmount       string `json:"vat_amount,omitempty"`
			Quantity        int    `json:"quantity,omitempty"`
			TotalPrice      string `json:"total_price,omitempty"`
			TotalWithVat    string `json:"total_with_vat,omitempty"`
		} `json:"products,omitempty"`
		VatRates []struct {
			Rate  float64 `json:"rate,omitempty"`
			Net   float64 `json:"net,omitempty"`
			Vat   float64 `json:"vat,omitempty"`
			Gross float64 `json:"gross,omitempty"`
		} `json:"vat_rates,omitempty"`
		Events    []any  `json:"events,omitempty"`
		ReceiptNo string `json:"receipt_no,omitempty"`
	} `json:"transaction_data,omitempty"`
	CardApplicationData struct {
		Name string `json:"name,omitempty"`
		Aid  string `json:"aid,omitempty"`
	} `json:"card_application_data,omitempty"`
	MerchantData struct {
		MerchantProfile struct {
			MerchantCode              string `json:"merchant_code,omitempty"`
			BusinessName              string `json:"business_name,omitempty"`
			CompanyRegistrationNumber string `json:"company_registration_number,omitempty"`
			Website                   string `json:"website,omitempty"`
			Email                     string `json:"email,omitempty"`
			Address                   struct {
				AddressLine1      string `json:"address_line1,omitempty"`
				AddressLine2      string `json:"address_line2,omitempty"`
				City              string `json:"city,omitempty"`
				Country           string `json:"country,omitempty"`
				CountryEnName     string `json:"country_en_name,omitempty"`
				CountryNativeName string `json:"country_native_name,omitempty"`
				PostCode          string `json:"post_code,omitempty"`
				Landline          string `json:"landline,omitempty"`
			} `json:"address,omitempty"`
		} `json:"merchant_profile,omitempty"`
		Locale string `json:"locale,omitempty"`
	} `json:"merchant_data,omitempty"`
	EmvData struct {
		Tvr    string `json:"tvr,omitempty"`
		Tsi    string `json:"tsi,omitempty"`
		Cvr    string `json:"cvr,omitempty"`
		Iad    string `json:"iad,omitempty"`
		Arc    string `json:"arc,omitempty"`
		Aid    string `json:"aid,omitempty"`
		Act    string `json:"act,omitempty"`
		Acv    string `json:"acv,omitempty"`
		Atc    string `json:"atc,omitempty"`
		Pan    string `json:"pan,omitempty"`
		TxType string `json:"tx_type,omitempty"`
	} `json:"emv_data,omitempty"`
	AcquirerData struct {
		Tid               string    `json:"tid,omitempty"`
		Mid               string    `json:"mid,omitempty"`
		AuthorizationCode string    `json:"authorization_code,omitempty"`
		Stan              int       `json:"stan,omitempty"`
		ReturnCode        string    `json:"return_code,omitempty"`
		LocalTime         time.Time `json:"local_time,omitempty"`
	} `json:"acquirer_data,omitempty"`
}
