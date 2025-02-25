package db

func (c *Connector) InitInternalDBIfNotExist() error {
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS `abra_sync_order_state` (`order_number` VARCHAR(255) PRIMARY KEY, `status` VARCHAR(255), `email_sent` BOOLEAN, `email` VARCHAR(255), `invoice_id` VARCHAR(255))")
	return err
}

func (c *Connector) QueryOrderProccesedState(orderNumber string) (string, error) {
	rows, err := c.db.Query("SELECT status FROM `abra_sync_order_state` WHERE order_number = ?", orderNumber)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var orderPaymentStatusID string
	for rows.Next() {
		err := rows.Scan(&orderPaymentStatusID)
		if err != nil {
			return "", err
		}
	}

	return orderPaymentStatusID, nil
}

func (c *Connector) UpdateOrderProccesedState(orderNumber string, status string) error {
	_, err := c.db.Exec("UPDATE `abra_sync_order_state` SET status = ? WHERE order_number = ?", status, orderNumber)
	return err
}

func (c *Connector) InsertOrderProccesedState(orderNumber string, status string, email string, invoiceID string) error {
	_, err := c.db.Exec("INSERT INTO `abra_sync_order_state` (order_number, status, email_sent, email, invoice_id) VALUES (?, ?, false, ?, ?)", orderNumber, status, email, invoiceID)
	return err
}

func (c *Connector) UpdateEmailSentState(orderNumber string) error {
	_, err := c.db.Exec("UPDATE `abra_sync_order_state` SET email_sent = ? WHERE order_number = ?", true, orderNumber)
	return err
}

func (c *Connector) InsertOrUpdateProcessedState(orderNumber string, status string, email string, invoiceID string) error {
	currentState, err := c.QueryOrderProccesedState(orderNumber)
	if err != nil {
		return err
	}
	if currentState == "" {
		return c.InsertOrderProccesedState(orderNumber, status, email, invoiceID)
	}
	return c.UpdateOrderProccesedState(orderNumber, status)
}

func (c *Connector) QueryUnsendOrders() ([]OrderInternalState, error) {
	rows, err := c.db.Query("SELECT order_number, status, email, invoice_id FROM `abra_sync_order_state` WHERE email_sent = false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []OrderInternalState{}
	for rows.Next() {
		var order OrderInternalState
		err := rows.Scan(&order.OrderNumber, &order.Status, &order.Email, &order.InvoiceID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
