package service

import (
	"errors"
	"gapstack-api/internal/model"
	"gapstack-api/internal/repository"
)

type TransactionService struct {
	Repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{Repo: repo}
}

func (s *TransactionService) Create(trx *model.Transaction) (*model.Transaction, error) {
	trx.Status = model.Status(model.StatusPending) // Additional business logic can be added here
	// e.g., validate currency, check sender/receiver accounts, etc.
	if trx.Amount <= 0 {
		return &model.Transaction{}, errors.New("amount must be positive amount")
	}
	err := s.Repo.CreateTransaction(trx)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

func (s *TransactionService) GetByID(id string) (*model.Transaction, error) {
	return s.Repo.GetTransactionByID(id)
}

func (s *TransactionService) List(limit, offset int) ([]model.Transaction, error) {
	return s.Repo.ListTransactions(limit, offset)
}

func (s *TransactionService) UpdateStatus(id string, newStatus model.Status) (*model.Transaction, error) {
	trx, err := s.Repo.GetTransactionByID(id)
	if err != nil {
		return nil, err
	}

	if trx.Status != model.StatusPending {
		return nil, errors.New("only pending transactions can be updated")
	}

	if newStatus != model.StatusCompleted && newStatus != model.StatusFailed {
		return nil, errors.New("invalid status transition")
	}

	// Update the database
	if err := s.Repo.UpdateTransactionStatus(id, model.Status(newStatus)); err != nil {
		return nil, err
	}

	// Reflect change in returned struct
	trx.Status = model.Status(newStatus)
	return trx, nil
}
