package db

import (
	"database/sql"
	"fmt"
)

type TX struct {
	*sql.Tx
}

func (tx *TX) SaveTransaction(height uint64, time string, transaction string) error {
	if tx.Tx == nil {
		panic(fmt.Errorf("SaveTransaction There is no transaction"))
	}
	_, err := tx.Exec("INSERT INTO transactions (height, time, transaction) VALUES ($1,$2,$3)", height, time, transaction)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("SaveTransaction failed to insert into transactions table: %v", err)
	}

	return nil
}

func (dbs *DBS) BeginTransaction() (*TX, error) {
	tx, err := dbs.Begin()
	if err != nil {
		return nil, fmt.Errorf("BeginTransaction failed to begin database transaction: %v", err)
	}

	return &TX{tx}, nil
}

func (tx *TX) Commit() error {
	err := tx.Tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Commit failed to commit database transaction: %v", err)
	}

	return nil
}

func (dbs *DBS) GetDone() (map[uint64]bool, error) {
	rows, err := dbs.Query("SELECT height FROM transactions")
	if err != nil {
		return nil, fmt.Errorf("GetDone Failed to query height from transactions: %v", err)
	}

	defer rows.Close()

	result := make(map[uint64]bool)
	var i uint64
	for rows.Next() {
		if err := rows.Scan(&i); err != nil {
			return nil, fmt.Errorf("GetDone Failed to scan rows: %v", err)
		}
		result[i] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetDone Failed to iterate rows: %v", err)
	}

	return result, nil
}
