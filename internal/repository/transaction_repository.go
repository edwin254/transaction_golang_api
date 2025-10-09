package repository

import (
    "gorm.io/gorm"
    "gapstack-api/internal/model"
)

type TransactionRepository struct {
    db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
    return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *model.Transaction) (*model.Transaction, error) {
    if err := r.db.Create(tx).Error; err != nil {
        return nil, err
    }
    return tx, nil
}

func (r *TransactionRepository) GetByID(id string) (*model.Transaction, error) {
    var tx model.Transaction
    if err := r.db.First(&tx, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &tx, nil
}

func (r *TransactionRepository) GetAll() ([]model.Transaction, error) {
    var txs []model.Transaction
    if err := r.db.Find(&txs).Error; err != nil {
        return nil, err
    }
    return txs, nil
}

func (r *TransactionRepository) Update(tx *model.Transaction) (*model.Transaction, error) {
    if err := r.db.Save(tx).Error; err != nil {
        return nil, err
    }
    return tx, nil
}
