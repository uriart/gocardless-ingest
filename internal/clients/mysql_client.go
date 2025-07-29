package clients

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocardless-ingest/internal/models"
)

type MySQLClient struct {
	db *sql.DB
}

// NewMySQLClient crea una nueva instancia y abre la conexión
func NewMySQLClient(dsn string) (*MySQLClient, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLClient{db: db}, nil
}

func (c *MySQLClient) TransactionExists(tx models.Transaction) (bool, error) {
	var exists bool
	err := c.db.QueryRow("SELECT EXISTS(SELECT 1 FROM transactions WHERE id = ?)", tx.ID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// InsertTransaction inserta una transacción si no existe
func (c *MySQLClient) InsertTransaction(tx models.Transaction) (sql.Result, error) {
	return c.db.Exec(`
		INSERT INTO transactions (id, booking_date, amount, currency, creditor_name, purpose_code, description, balance_after)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		tx.ID, tx.BookingDate, tx.TransactionAmount.Amount,
		tx.TransactionAmount.Currency, tx.CreditorName,
		tx.PurposeCode, tx.Description, tx.BalanceAfter.BalanceAmount.Amount,
	)
}

// Close cierra la conexión a la base de datos
func (c *MySQLClient) Close() error {
	return c.db.Close()
}
