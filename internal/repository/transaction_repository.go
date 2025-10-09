package repository

import (
	"database/sql"
	"fmt"

	"yourapp/internal/model"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

// Create
func (r *TransactionRepository) Create(tx *model.Transaction) error {
	tx.ID = uuid.New().String()
	query := `INSERT INTO transactions (id, amount, currency, sender, receiver, status)
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, tx.ID, tx.Amount, tx.Currency, tx.Sender, tx.Receiver, tx.Status)
	return err
}

// Retrieve by ID
func (r *TransactionRepository) GetByID(id string) (*model.Transaction, error) {
	query := `SELECT id, amount, currency, sender, receiver, status FROM transactions WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var tx model.Transaction
	if err := row.Scan(&tx.ID, &tx.Amount, &tx.Currency, &tx.Sender, &tx.Receiver, &tx.Status); err != nil {
		return nil, err
	}
	return &tx, nil
}

// List with pagination
func (r *TransactionRepository) List(limit, offset int) ([]model.Transaction, error) {
	query := `SELECT id, amount, currency, sender, receiver, status FROM transactions
	          ORDER BY ROWID DESC LIMIT ? OFFSET ?`
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.ID, &tx.Amount, &tx.Currency, &tx.Sender, &tx.Receiver, &tx.Status); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// Update status (only valid transitions)
func (r *TransactionRepository) UpdateStatus(id string, newStatus model.Status) error {
	query := `UPDATE transactions SET status = ? 
	          WHERE id = ? AND status = 'pending'`
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
