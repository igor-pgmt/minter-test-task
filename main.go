package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/igor-pgmt/minter-test-task/client"
	"github.com/igor-pgmt/minter-test-task/db"
	"github.com/igor-pgmt/minter-test-task/server"
	"github.com/igor-pgmt/minter-test-task/utils"
)

var domainName = "minter-node-1.mainnet.minter.network"
var baseURL = "https://" + domainName + "/block?height="
var fastHTTPclientURL = domainName + ":80"

var poolCap uint64
var minBlock uint64
var maxBlock uint64
var responseSize uint64

var dbHost string
var dbPort string
var dbName string
var dbTableName string
var dbUser string
var dbPass string
var clean bool
var continueParsing bool
var skipParsing bool

var parsed uint64

func init() {

	var err error

	domainName = os.Getenv("domainName")
	if len(domainName) == 0 {
		panic(fmt.Errorf("no domainName env set"))
	}

	poolCapString := os.Getenv("poolCap")
	if len(poolCapString) == 0 {
		panic(fmt.Errorf("no poolCap env set"))
	}
	poolCap, err = strconv.ParseUint(poolCapString, 10, 64)
	if err != nil {
		panic(fmt.Errorf("incorrect poolCap env set"))
	}
	if poolCap < 1 {
		panic(fmt.Errorf("poolCap must be greater then 0"))
	}

	minBlockString := os.Getenv("minBlock")
	if len(minBlockString) == 0 {
		panic(fmt.Errorf("no minBlock env set"))
	}
	minBlock, err = strconv.ParseUint(minBlockString, 10, 64)
	if err != nil {
		panic(fmt.Errorf("incorrect minBlock env set"))
	}
	if minBlock < 1 {
		panic(fmt.Errorf("minBlock must be greater then 0"))
	}

	maxBlockString := os.Getenv("maxBlock")
	if len(maxBlockString) == 0 {
		panic(fmt.Errorf("no maxBlock env set"))
	}
	maxBlock, err = strconv.ParseUint(maxBlockString, 10, 64)
	if err != nil {
		panic(fmt.Errorf("incorrect maxBlock env set"))
	}
	if maxBlock < 2 {
		panic(fmt.Errorf("maxBlock must be greater then 1"))
	}

	responseSizeString := os.Getenv("responseSize")
	if len(responseSizeString) == 0 {
		panic(fmt.Errorf("no responseSize env set"))
	}
	responseSize, err = strconv.ParseUint(responseSizeString, 10, 64)
	if err != nil {
		panic(fmt.Errorf("incorrect responseSize env set"))
	}

	dbHost = os.Getenv("dbHost")
	if len(dbHost) == 0 {
		panic(fmt.Errorf("no dbHost env set"))
	}
	dbPort = os.Getenv("dbPort")
	if len(dbPort) == 0 {
		panic(fmt.Errorf("no dbPort env set"))
	}
	dbName = os.Getenv("dbName")
	if len(dbName) == 0 {
		panic(fmt.Errorf("no dbName env set"))
	}
	dbUser = os.Getenv("dbUser")
	if len(dbUser) == 0 {
		panic(fmt.Errorf("no dbUser env set"))
	}
	dbPass = os.Getenv("dbPass")
	if len(dbPass) == 0 {
		panic(fmt.Errorf("no dbPass env set"))
	}
	cleanString := os.Getenv("clean")
	if len(cleanString) == 0 {
		panic(fmt.Errorf("no clean env set"))
	}
	clean, err = strconv.ParseBool(cleanString)
	if err != nil {
		panic(fmt.Errorf("incorrect clean env set"))
	}
	continueParsingString := os.Getenv("continueParsing")
	if len(continueParsingString) == 0 {
		panic(fmt.Errorf("no continueParsing env set"))
	}
	continueParsing, err = strconv.ParseBool(continueParsingString)
	if err != nil {
		panic(fmt.Errorf("incorrect continueParsing env set"))
	}
	skipParsingString := os.Getenv("skipParsing")
	if len(skipParsingString) == 0 {
		panic(fmt.Errorf("no skipParsing env set"))
	}
	skipParsing, err = strconv.ParseBool(skipParsingString)
	if err != nil {
		panic(fmt.Errorf("incorrect skipParsing env set"))
	}

}

func main() {

	rLimit, err := utils.MaxOpenFiles()
	if err != nil {
		panic(fmt.Errorf("Failed to set MaxOpenFiles: %v", err))
	}

	if poolCap > rLimit.Cur {
		poolCap = rLimit.Cur - rLimit.Cur/20
	}

	timeRun := time.Now()

	db1, err := db.DBConnect(dbHost, dbPort, dbUser, dbPass, dbName, int(poolCap))
	if err != nil {
		panic(fmt.Errorf("Failed to connect to the database %s: %v", dbName, err))
	}

	exists, err := db1.TableExists()
	if err != nil {
		panic(fmt.Errorf("Failed to check if table exists: %v", err))
	}

	if !exists {
		err := db1.CreateTable()
		if err != nil {
			panic(fmt.Errorf("Failed to create table: %v", err))
		}
	}

	if clean {
		err = db1.Clean()
		if err != nil {
			panic(fmt.Errorf("Failed to clean table: %v", err))
		}
	}

	if !skipParsing {

		var previous map[uint64]bool
		if continueParsing {
			previous, err = db1.GetDone()
			if err != nil {
				panic(fmt.Errorf("Failed to get previously parsed block heights: %v", err))
			}
		}

		c := client.NewClient(fastHTTPclientURL, responseSize, poolCap)

		done := make(chan struct{})
		go func(parsed *uint64, done <-chan struct{}) {
			for {
				select {
				case <-done:
					return
				default:
					log.Printf("done: %g%%\n", float64(atomic.LoadUint64(parsed))*float64(100)/float64(maxBlock))
				}
				time.Sleep(time.Second)
			}
		}(&parsed, done)

		wg := &sync.WaitGroup{}
		var url = baseURL
		for blockNum := minBlock; blockNum <= maxBlock; blockNum++ {
			if previous[blockNum] {
				atomic.AddUint64(&parsed, 1)
				continue
			}

			url = baseURL + fmt.Sprint(blockNum)
			wg.Add(1)

			go func(wg *sync.WaitGroup, c *client.FastHTTPClient, blockNum uint64, db *db.DBS, url string) {
				r := &client.Response{}
				var err error
				for r, err = c.GetBlock(url, blockNum); err != nil; r, err = c.GetBlock(url, blockNum) {
					// suppressing usual timeout errors
					if !strings.Contains(err.Error(), "dialing to the given TCP address timed out") && !strings.Contains(err.Error(), "Unexpected status code: 502. Expecting 200") {
						log.Printf("block: %d, parsed:%v ERR: %v", blockNum, atomic.LoadUint64(&parsed), err)
					}
				}

				tx, err := db.BeginTransaction()
				if err != nil {
					panic(fmt.Errorf("Falied to begin transaction for block %d: %v", blockNum, err))
				}

				for _, transaction := range r.Result.Transactions {
					if transaction.Type == 1 || transaction.Type == 3 {

						t, err := json.Marshal(transaction)
						if err != nil {
							log.Printf("masrsal t err: %v\n", err)
						}

						err = tx.SaveTransaction(blockNum, r.Result.Time, string(t))
						if err != nil {
							panic(fmt.Errorf("Falied to save transaction for block %d: %v", blockNum, err))
						}
					}
				}

				err = tx.Commit()
				if err != nil {
					panic(fmt.Errorf("Falied to commit transactions for block %d: %v", blockNum, err))
				}

				atomic.AddUint64(&parsed, 1)
				wg.Done()
			}(wg, c, blockNum, db1, url)
		}
		wg.Wait()

		done <- struct{}{}
		close(done)

		log.Printf("done: %g%%\n", float64(atomic.LoadUint64(&parsed))*float64(100)/float64(maxBlock))
		log.Println("total parsing time: ", time.Since(timeRun))

	}

	log.Println("Starting server ...")

	db2, err := db.DBConnect(dbHost, dbPort, dbUser, dbPass, dbName, int(poolCap))
	if err != nil {
		panic(err)
	}
	server := server.NewServer(db2)
	server.Run()
}
