package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBS struct {
	*sql.DB
}

func DBConnect(host, port, user, pass, dbname string, maxIdleConn int) (*DBS, error) {

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable binary_parameters=yes",
		host, port, user, pass, dbname)
	dbc, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to open postgres connection: %q", err)
	}

	err = dbc.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to check postgres connection: %q", err)
	}

	dbc.SetMaxOpenConns(0)
	dbc.SetMaxIdleConns(maxIdleConn)

	return &DBS{dbc}, nil
}

func (dbs *DBS) CreateTable() error {

	tx, err := dbs.Begin()
	if err != nil {
		return fmt.Errorf("CreateTable Failed to begin transaction: %v", err)
	}

	_, err = tx.Exec("CREATE TABLE transactions (height INTEGER NOT NULL, time timestamp with time zone NOT NULL, transaction jsonb NOT NULL);")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("CreateTable Failed to create table transactions: %v", err)
	}

	_, err = tx.Exec("CREATE INDEX on transactions USING btree ((transaction ->> 'hash'));")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("CreateTable Failed to create index on table transactions: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("CreateTable Failed to commit transaction: %v", err)
	}

	return nil
}

func (dbs *DBS) TableExists() (bool, error) {
	var exists bool
	err := dbs.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'transactions');").Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("GetDone Failed to iterate rows: %v", err)
	}

	return exists, nil
}

func (dbs *DBS) Clean() error {
	_, err := dbs.Exec("TRUNCATE TABLE transactions;")
	if err != nil {
		return fmt.Errorf("Clean Failed to truncate table transactions: %v", err)
	}

	return nil
}
