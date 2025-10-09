package seeder

import (
	"database/sql"
	"log"

	"gapstack-api/internal/model"

	"github.com/google/uuid"
)

func SeedTransactions(db *sql.DB) error {
	// Check if DB already has data
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("✅ Seeder skipped — transactions table already populated")
		return nil
	}

	// Sample data to insert
	sampleData := []model.Transaction{
		{ID: uuid.New().String(), Amount: 50.00, Currency: "USD", Sender: "alice", Receiver: "bob", Status: model.StatusPending},
		{ID: uuid.New().String(), Amount: 100.50, Currency: "EUR", Sender: "carol", Receiver: "dave", Status: model.StatusCompleted},
		{ID: uuid.New().String(), Amount: 230.00, Currency: "KES", Sender: "eve", Receiver: "frank", Status: model.StatusFailed},
		{ID: uuid.New().String(), Amount: 70.00, Currency: "USD", Sender: "alice", Receiver: "george", Status: model.StatusCompleted},
		{ID: uuid.New().String(), Amount: 90.00, Currency: "USD", Sender: "bob", Receiver: "alice", Status: model.StatusPending},
	}

	stmt, err := db.Prepare(`
		INSERT INTO transactions (id, amount, currency, sender, receiver, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, trx := range sampleData {
		_, err := stmt.Exec(trx.ID, trx.Amount, trx.Currency, trx.Sender, trx.Receiver, trx.Status)
		if err != nil {
			return err
		}
	}

	log.Printf("✅ Seeder completed — inserted %d transactions\n", len(sampleData))
	return nil
}
