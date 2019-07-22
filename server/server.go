package server

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/igor-pgmt/minter-test-task/db"
)

type srv struct {
	db *db.DBS
}

func NewServer(db *db.DBS) *srv {
	return &srv{db}
}

func (s *srv) Run() {

	http.HandleFunc("/api/transactions/from/", s.GetTransactions)
	http.HandleFunc("/api/transactions/get", s.GetByTime)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Errorf("Failed to start http server: %v", err))
	}
}

func (s *srv) GetByTime(w http.ResponseWriter, r *http.Request) {
	time1str := r.URL.Query().Get("time1")
	time2str := r.URL.Query().Get("time2")
	if time1str == "" || time2str == "" {
		http.Error(w, "You need to set time1 and time2", http.StatusBadRequest)
		return
	}

	time1, err := time.Parse("2006-01-02 15:04:05.000000-07", time1str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	time2, err := time.Parse("2006-01-02 15:04:05.000000-07", time2str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions, err := s.db.GetByTime(time1, time2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ResponseTime{}

	for from, tsFrom := range transactions {
		ts := Transactions{From: from, Sum: new(big.Int)}
		for _, transaction := range tsFrom {
			ts.Transactions = append(ts.Transactions, transaction)
			if transaction.Type == 1 {
				i, _ := new(big.Int).SetString(transaction.Data.Value, 10)
				ts.Sum = ts.Sum.Add(ts.Sum, i)
			}
		}
		res.Result.Transactions = append(res.Result.Transactions, ts)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *srv) GetTransactions(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	from := path[len(path)-1]
	if from == "" {
		http.Error(w, "You need to set FROM", http.StatusBadRequest)
		return
	}

	transactions, err := s.db.GetTransactions(from)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sum := new(big.Int)
	for _, t := range transactions {
		if t.Type == 1 {
			i, _ := new(big.Int).SetString(t.Data.Value, 10)
			sum = sum.Add(sum, i)
		}
	}

	res := Response{}
	res.Result.Sum = sum
	res.Result.Transactions = transactions

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
