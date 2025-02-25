package db

type Order struct {
	OrderID               int     `json:"order_id,omitempty"`
	ReservationID         int     `json:"reservation_id,omitempty"`
	ID                    int     `json:"id,omitempty"`
	LocaleID              int     `json:"locale_id,omitempty"`
	UserID                int     `json:"user_id,omitempty"`
	OrderNumber           string  `json:"order_number,omitempty"`
	Created               string  `json:"created,omitempty"`
	PaymentType           string  `json:"payment_type,omitempty"`
	TotalPrice            float64 `json:"total_price,omitempty"`
	CsobGwID              string  `json:"csob_gw_id,omitempty"`
	PaymentPrice          float64 `json:"payment_price,omitempty"`
	PaymentVat            int     `json:"payment_vat,omitempty"`
	OrderPaymentStatusID  int     `json:"order_payment_status_id,omitempty"`
	PaymentSettings       string  `json:"payment_settings,omitempty"`
	CurrencyID            int     `json:"currency_id,omitempty"`
	InvoiceNum            string  `json:"invoice_num,omitempty"`
	PaymentReceivedAt     string  `json:"payment_received_at,omitempty"`
	InvoiceCreated        string  `json:"invoice_created,omitempty"`
	CreditNoteNum         any     `json:"credit_note_num,omitempty"`
	CreditNoteCreated     any     `json:"credit_note_created,omitempty"`
	OrderCreatedEmailSend int     `json:"order_created_email_send,omitempty"`
	ServiceID             int     `json:"service_id,omitempty"`
	Start                 string  `json:"start,omitempty"`
	End                   string  `json:"end,omitempty"`
	Options               string  `json:"options,omitempty"`
	CanceledByID          any     `json:"canceled_by_id,omitempty"`
	Canceled              int     `json:"canceled,omitempty"`
	CancelReason          any     `json:"cancel_reason,omitempty"`
	PlaceID               int     `json:"place_id,omitempty"`
	Note                  any     `json:"note,omitempty"`
	Vat                   int     `json:"vat,omitempty"`
	Price                 float64 `json:"price,omitempty"`
	GoogleID              any     `json:"google_id,omitempty"`
	Capacity              int     `json:"capacity,omitempty"`
}

type User struct {
	ID                  int    `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Email               string `json:"email,omitempty"`
	Password            string `json:"password,omitempty"`
	Position            any    `json:"position,omitempty"`
	Phone               string `json:"phone,omitempty"`
	Admin               int    `json:"admin,omitempty"`
	Token               any    `json:"token,omitempty"`
	TokenExpiration     any    `json:"token_expiration,omitempty"`
	DateOfBirth         any    `json:"date_of_birth,omitempty"`
	Note                any    `json:"note,omitempty"`
	Surname             string `json:"surname,omitempty"`
	Image               any    `json:"image,omitempty"`
	Tin                 any    `json:"tin,omitempty"`
	VatID               any    `json:"vat_id,omitempty"`
	NotificationTimeout int    `json:"notification_timeout,omitempty"`
	UnsubscribeEmail    int    `json:"unsubscribe_email,omitempty"`
	UnsubscribeSms      int    `json:"unsubscribe_sms,omitempty"`
	NotificationAllowed int    `json:"notification_allowed,omitempty"`
	LastPhoneContact    any    `json:"last_phone_contact,omitempty"`
	Vip                 int    `json:"vip,omitempty"`
	SaunaBoda           int    `json:"sauna_boda,omitempty"`
	AddressID           any    `json:"address_id,omitempty"`
	AddressUserID       int    `json:"address_user_id,omitempty"`
	Company             any    `json:"company,omitempty"`
	AddressName         string `json:"address_name,omitempty"`
	Street              string `json:"street,omitempty"`
	HouseNumber         string `json:"house_number,omitempty"`
	City                string `json:"city,omitempty"`
	ZipCode             string `json:"zip_code,omitempty"`
	Region              any    `json:"region,omitempty"`
	Country             string `json:"country,omitempty"`
	Type                any    `json:"type,omitempty"`
}

type Ticket struct {
	ID                      int     `json:"id,omitempty"`
	TicketID                int     `json:"ticket_id,omitempty"`
	Price                   float64 `json:"price,omitempty"`
	OrderID                 int     `json:"order_id,omitempty"`
	Vat                     int     `json:"vat,omitempty"`
	IsRenew                 int     `json:"is_renew,omitempty"`
	ValidityAlreadyExtended int     `json:"validity_already_extended,omitempty"`
	LocaleID                int     `json:"locale_id,omitempty"`
	UserID                  int     `json:"user_id,omitempty"`
	OrderNumber             string  `json:"order_number,omitempty"`
	Created                 string  `json:"created,omitempty"`
	PaymentType             string  `json:"payment_type,omitempty"`
	TotalPrice              float64 `json:"total_price,omitempty"`
	CsobGwID                string  `json:"csob_gw_id,omitempty"`
	PaymentPrice            float64 `json:"payment_price,omitempty"`
	PaymentVat              int     `json:"payment_vat,omitempty"`
	OrderPaymentStatusID    int     `json:"order_payment_status_id,omitempty"`
	PaymentSettings         string  `json:"payment_settings,omitempty"`
	CurrencyID              int     `json:"currency_id,omitempty"`
	InvoiceNum              string  `json:"invoice_num,omitempty"`
	PaymentReceivedAt       string  `json:"payment_received_at,omitempty"`
	InvoiceCreated          string  `json:"invoice_created,omitempty"`
	CreditNoteNum           any     `json:"credit_note_num,omitempty"`
	CreditNoteCreated       any     `json:"credit_note_created,omitempty"`
	OrderCreatedEmailSend   int     `json:"order_created_email_send,omitempty"`
}

type OrderInternalState struct {
	OrderNumber string `json:"order_number"`
	Status      string `json:"status"`
	EmailSent   int    `json:"email_sent"`
	Email       string `json:"email"`
	InvoiceID   string `json:"invoice_id"`
}
