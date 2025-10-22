package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"gapstack-api/internal/model"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

// validateTransaction checks the transaction fields for validity.
func validateTransaction(trx *model.Transaction) error {
	if trx.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if trx.Currency == "" {
		return errors.New("currency is required")
	}
	if trx.Sender == "" {
		return errors.New("sender is required")
	}
	if trx.Receiver == "" {
		return errors.New("receiver is required")
	}
	if trx.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

// Create
func (r *TransactionRepository) CreateTransaction(trx *model.Transaction) error {
	// Input validation
	if err := validateTransaction(trx); err != nil {
		return fmt.Errorf("invalid transaction: %w", err)
	}

	// Generate UUID if not provided
	if trx.ID == "" {
		trx.ID = uuid.New().String()
	}

	query := `INSERT INTO transactions (id, amount, currency, sender, receiver, status)
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, trx.ID, trx.Amount, trx.Currency, trx.Sender, trx.Receiver, trx.Status)

	if err != nil {
		log.Fatalf("failed to insert transaction: %w", err)
		return err // Return nil transaction on error
	}
	return nil
}

// Retrieve by ID
func (r *TransactionRepository) GetTransactionByID(id string) (*model.Transaction, error) {
	query := `SELECT id, amount, currency, sender, receiver, status FROM transactions WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var trx model.Transaction
	if err := row.Scan(&trx.ID, &trx.Amount, &trx.Currency, &trx.Sender, &trx.Receiver, &trx.Status); err != nil {
		return nil, err
	}
	return &trx, nil
}

// List with pagination
func (r *TransactionRepository) ListTransactions(limit, offset int) ([]model.Transaction, error) {
	query := `SELECT * FROM transactions
	          ORDER BY id DESC LIMIT $1 OFFSET $2`
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var trx model.Transaction
		if err := rows.Scan(&trx.ID, &trx.Amount, &trx.Currency, &trx.Sender, &trx.Receiver, &trx.Status, &trx.CreatedAt, &trx.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, trx)
	}
	return transactions, nil
}

// Update status (only valid transitions)
func (r *TransactionRepository) UpdateTransactionStatus(id string, newStatus model.Status) error {
	query := `UPDATE transactions SET status = $1 
	          WHERE id = $2 AND status = 'pending'`

	result, err := r.DB.Exec(query, newStatus, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invalid status transition or transaction not found")
	}
	return nil
}

// Total amount of completed transactions per user
func (r *TransactionRepository) TotalCompletedPerUser() (map[string]float64, error) {
	query := `SELECT sender, SUM(amount) as total FROM transactions 
	          WHERE status = 'completed' GROUP BY sender`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := make(map[string]float64)
	for rows.Next() {
		var user string
		var total float64
		if err := rows.Scan(&user, &total); err != nil {
			return nil, err
		}
		totals[user] = total
	}
	return totals, nil
}

// Top 5 users by transaction volume (last 30 days)
func (r *TransactionRepository) TopUsersLast30Days() ([]string, error) {
	query := `SELECT sender, SUM(amount) as total FROM transactions
	          WHERE status = 'completed' AND datetime(created_at) >= datetime('now', '-30 days')
	          GROUP BY sender ORDER BY total DESC LIMIT 5`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var sender string
		var total float64
		if err := rows.Scan(&sender, &total); err != nil {
			return nil, err
		}
		users = append(users, sender)
	}
	return users, nil
}

// Users with more than 3 failed transactions
func (r *TransactionRepository) UsersWithMoreThan3Failed() ([]string, error) {
	query := `SELECT sender FROM transactions
	          WHERE status = 'failed'
	          GROUP BY sender
	          HAVING COUNT(*) > 3`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var sender string
		if err := rows.Scan(&sender); err != nil {
			return nil, err
		}
		users = append(users, sender)
	}
	return users, nil
}
