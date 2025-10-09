package config

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gapstack-api/internal/model"
)

func InitDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("transactions.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    db.AutoMigrate(&model.Transaction{})
    return db, nil
}
