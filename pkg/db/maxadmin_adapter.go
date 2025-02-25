package db

func (c *Connector) QueryPayedOrdersInYear(date string) ([]Order, error) {
	rows, err := c.db.Query("SELECT o.id AS order_id, r.id AS reservation_id, o.locale_id, o.user_id, o.order_number, o.created, o.payment_type, o.total_price, o.csob_gw_id, o.payment_price, o.payment_vat, o.order_payment_status_id, o.payment_settings, o.currency_id, o.invoice_num, o.payment_received_at, o.invoice_created, o.credit_note_num, o.credit_note_created, o.order_created_email_send, r.service_id, r.start, r.end, r.options, r.canceled_by_id, r.canceled, r.cancel_reason, r.place_id, r.note, r.vat, r.price, r.google_id, r.capacity FROM `order` o RIGHT JOIN reservation r ON o.id = r.order_id WHERE o.payment_type = 'csob-gateway' AND o.created >= STR_TO_DATE('" + date + "', '%Y-%m-%d') AND o.payment_received_at IS NOT NULL AND r.end < NOW() AND r.canceled = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.OrderID, &order.ReservationID, &order.LocaleID, &order.UserID, &order.OrderNumber, &order.Created, &order.PaymentType, &order.TotalPrice, &order.CsobGwID, &order.PaymentPrice, &order.PaymentVat, &order.OrderPaymentStatusID, &order.PaymentSettings, &order.CurrencyID, &order.InvoiceNum, &order.PaymentReceivedAt, &order.InvoiceCreated, &order.CreditNoteNum, &order.CreditNoteCreated, &order.OrderCreatedEmailSend, &order.ServiceID, &order.Start, &order.End, &order.Options, &order.CanceledByID, &order.Canceled, &order.CancelReason, &order.PlaceID, &order.Note, &order.Vat, &order.Price, &order.GoogleID, &order.Capacity)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (c *Connector) QueryPaysChipsInYear(date string) ([]ChipOrder, error) {
	rows, err := c.db.Query("SELECT o.id AS order_id, o.locale_id, o.user_id, o.order_number, o.created, o.payment_type, o.total_price, o.csob_gw_id, o.payment_price, o.payment_vat, o.order_payment_status_id, o.payment_settings, o.currency_id, o.invoice_num, o.payment_received_at, o.invoice_created, o.credit_note_num, o.credit_note_created, o.order_created_email_send, coi.chip_product_id, coi.price, coi.vat, coi.was_notified, coi.admins_notified, coi.can_be_dispensed, coi.dispensed_at, coi.was_notified_dispenser FROM `order` o INNER JOIN chip_order_item coi ON o.id = coi.order_id WHERE o.payment_type = 'csob-gateway' AND o.created >= STR_TO_DATE('" + date + "', '%Y-%m-%d') AND o.payment_received_at IS NOT NULL AND coi.price != 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chips := []ChipOrder{}
	for rows.Next() {
		var chip ChipOrder
		err := rows.Scan(&chip.OrderID, &chip.LocaleID, &chip.UserID, &chip.OrderNumber, &chip.Created, &chip.PaymentType, &chip.TotalPrice, &chip.CsobGWID, &chip.PaymentPrice, &chip.PaymentVAT, &chip.OrderPaymentStatusID, &chip.PaymentSettings, &chip.CurrencyID, &chip.InvoiceNum, &chip.PaymentReceivedAt, &chip.InvoiceCreated, &chip.CreditNoteNum, &chip.CreditNoteCreated, &chip.OrderCreatedEmailSend, &chip.ChipProductID, &chip.Price, &chip.VAT, &chip.WasNotified, &chip.AdminsNotified, &chip.CanBeDispensed, &chip.DispensedAt, &chip.WasNotifiedDispenser)
		if err != nil {
			return nil, err
		}
		chips = append(chips, chip)
	}

	return chips, nil
}

func (c *Connector) QueryPayedTicketsInYear(date string) ([]Ticket, error) {
	rows, err := c.db.Query("SELECT t.id AS ticket_id, t.price, t.order_id, t.vat, t.is_renew, t.validity_already_extended, o.locale_id, o.user_id, o.order_number, o.created, o.payment_type, o.total_price, o.csob_gw_id, o.payment_price, o.payment_vat, o.order_payment_status_id, o.payment_settings, o.currency_id, o.invoice_num, o.payment_received_at, o.invoice_created, o.credit_note_num, o.credit_note_created, o.order_created_email_send FROM ticket_order_item t RIGHT JOIN `order` o ON o.id = t.order_id WHERE o.payment_type = 'csob-gateway' AND o.order_payment_status_id = 1 AND o.created >= STR_TO_DATE('" + date + "', '%Y-%m-%d') AND t.id IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := []Ticket{}
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(&ticket.ID, &ticket.Price, &ticket.OrderID, &ticket.Vat, &ticket.IsRenew, &ticket.ValidityAlreadyExtended, &ticket.LocaleID, &ticket.UserID, &ticket.OrderNumber, &ticket.Created, &ticket.PaymentType, &ticket.TotalPrice, &ticket.CsobGwID, &ticket.PaymentPrice, &ticket.PaymentVat, &ticket.OrderPaymentStatusID, &ticket.PaymentSettings, &ticket.CurrencyID, &ticket.InvoiceNum, &ticket.PaymentReceivedAt, &ticket.InvoiceCreated, &ticket.CreditNoteNum, &ticket.CreditNoteCreated, &ticket.OrderCreatedEmailSend)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (c *Connector) QueryUserByID(id int) (*User, error) {
	rows, err := c.db.Query("SELECT u.id AS user_id, u.name AS user_name, u.email, u.password, u.position, u.phone, u.admin, u.token, u.token_expiration, u.date_of_birth, u.note, u.surname, u.image, u.tin, u.vat_id, u.notification_timeout, u.unsubscribe_email, u.unsubscribe_sms, u.notification_allowed, u.last_phone_contact, u.vip, u.sauna_boda, ua.id AS address_id, ua.user_id AS address_user_id, ua.company, ua.name AS address_name, ua.street, ua.house_number, ua.city, ua.zip_code, ua.region, ua.country, ua.type FROM user u LEFT JOIN user_address ua ON u.id = ua.user_id WHERE u.id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &User{}
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Position, &user.Phone, &user.Admin, &user.Token, &user.TokenExpiration, &user.DateOfBirth, &user.Note, &user.Surname, &user.Image, &user.Tin, &user.VatID, &user.NotificationTimeout, &user.UnsubscribeEmail, &user.UnsubscribeSms, &user.NotificationAllowed, &user.LastPhoneContact, &user.Vip, &user.SaunaBoda, &user.AddressID, &user.AddressUserID, &user.Company, &user.AddressName, &user.Street, &user.HouseNumber, &user.City, &user.ZipCode, &user.Region, &user.Country, &user.Type)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
