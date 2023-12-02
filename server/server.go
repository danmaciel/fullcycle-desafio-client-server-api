package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danmaciel/client-server-api/helper"
	"github.com/danmaciel/client-server-api/model"
	_ "github.com/mattn/go-sqlite3"
)

// nos testes tive que colocar pra 600ms, pois o timeout estava estourando
const timeoutForApiRequest = time.Millisecond * 600
const timeoutForDatabaseIsertation = time.Millisecond * 10
const requestUri = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {

	initDb()

	log.Println("Server escutando em localhost:8080")
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(res http.ResponseWriter, req *http.Request) {

	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

	var bid model.Bid

	err := json.Unmarshal(helper.GetResponseApi(requestUri, timeoutForApiRequest), &bid)

	if err != nil {
		log.Fatal(err)
	}

	result := bid.USDBRL.Bid

	fmt.Printf("\nValor %v\n", result)

	db := getDb()

	defer db.Close()

	err = insertBidInDb(db, result)

	if err != nil {
		log.Fatal(err)
	}

	apiResponse(res, result)
}

func apiResponse(res http.ResponseWriter, result string) {
	valueMap := make(map[string]string)

	valueMap["bid"] = result

	re, err := json.Marshal(valueMap)
	if err != nil {
		log.Fatal(err)
	}

	res.Write(re)
}

func getDb() *sql.DB {

	db, err := sql.Open("sqlite3", "./goexpert.db")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func insertBidInDb(db *sql.DB, value string) error {
	ctx, cancelFunc := helper.GetCtxWithTimout(timeoutForDatabaseIsertation)
	defer cancelFunc()

	stmt, err := db.Prepare("Insert into cotacao (bid) values(?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, value)

	return err
}

func initDb() {
	database := getDb()
	ctx, cancelFunc := helper.GetCtxWithTimout(timeoutForDatabaseIsertation)

	defer cancelFunc()
	defer database.Close()

	_, err := database.ExecContext(ctx, "create table if not exists cotacao(bid varchar(10))")
	if err != nil {
		log.Fatal(err)
	}
}
