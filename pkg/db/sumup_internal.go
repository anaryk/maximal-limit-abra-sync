package db

import (
	"database/sql"
)

func (c *Connector) InitSumupTransactionStateTable() error {
	_, err := c.db.Exec(`
		CREATE TABLE IF NOT EXISTS sumup_transaction_state (
			transaction_id VARCHAR(255) PRIMARY KEY,
			imported BOOLEAN
		)`)
	return err
}

func (c *Connector) IsTransactionImported(transactionID string) (bool, error) {
	var imported bool
	err := c.db.QueryRow(`SELECT imported FROM sumup_transaction_state WHERE transaction_id = ?`, transactionID).Scan(&imported)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return imported, nil
}

func (c *Connector) InsertOrUpdateTransactionState(transactionID string, imported bool) error {
	_, err := c.db.Exec(`
		INSERT INTO sumup_transaction_state (transaction_id, imported)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE imported = VALUES(imported)
	`, transactionID, imported)
	return err
}
