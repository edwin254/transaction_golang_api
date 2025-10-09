package service

import (
    "errors"
    "gapstack-api/internal/model"
    "gapstack-api/internal/repository"
)

type TransactionService struct {
    repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
    return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(tx *model.Transaction) (*model.Transaction, error) {
    tx.Status = "pending"
    return s.repo.Create(tx)
}

func (s *TransactionService) GetTransactionByID(id string) (*model.Transaction, error) {
    return s.repo.GetByID(id)
}

func (s *TransactionService) ListTransactions() ([]model.Transaction, error) {
    return s.repo.GetAll()
}

func (s *TransactionService) UpdateTransactionStatus(id string, status string) (*model.Transaction, error) {
    tx, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }

    if tx.Status == "pending" && (status == "completed" || status == "failed") {
        tx.Status = status
        return s.repo.Update(tx)
    }
    return nil, errors.New("invalid status transition")
}
