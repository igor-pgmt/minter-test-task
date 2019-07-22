package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/igor-pgmt/minter-test-task/client"
)

func (dbs *DBS) GetByTime(time1, time2 time.Time) (map[string][]client.Transaction, error) {

	rows, err := dbs.Query("SELECT transaction FROM transactions WHERE time BETWEEN $1 AND $2 ORDER BY time", time1, time2)
	if err != nil {
		return nil, fmt.Errorf("GetByTime Failed to get transactions from the database between %v and %v: %v", time1, time2, err)
	}

	defer rows.Close()

	res := make(map[string][]client.Transaction)

	var ts string
	for rows.Next() {
		err := rows.Scan(&ts)
		if err != nil {
			return nil, fmt.Errorf("GetByTime Failed to scan transaction into the string: %v", err)
		}

		transaction := client.Transaction{}

		err = json.Unmarshal([]byte(ts), &transaction)
		if err != nil {
			return nil, fmt.Errorf("GetByTime Failed to unmarshal transaction: %v", err)
		}
		res[transaction.From] = append(res[transaction.From], transaction)
	}

	return res, nil
}

func (dbs *DBS) GetTransactions(from string) ([]client.Transaction, error) {

	rows, err := dbs.Query("SELECT transaction FROM transactions WHERE transaction->>'from' = $1;", from)
	if err != nil {
		return nil, fmt.Errorf("GetTransactions Failed to get transactions from the database from %s: %v", from, err)
	}

	defer rows.Close()

	var transactions []client.Transaction
	var ts string
	for rows.Next() {
		err := rows.Scan(&ts)
		if err != nil {
			return nil, fmt.Errorf("GetTransactions Failed to scan transaction into the string: %v", err)
		}

		transaction := client.Transaction{}

		err = json.Unmarshal([]byte(ts), &transaction)
		if err != nil {
			return nil, fmt.Errorf("GetTransactions Failed to unmarshal transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
