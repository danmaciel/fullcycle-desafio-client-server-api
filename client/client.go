package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danmaciel/client-api/helper"
	"github.com/danmaciel/client-api/model"
)

func main() {
	timeout := time.Millisecond * 1000
	requestUri := "http://localhost:8080/cotacao"

	var bid model.Bid

	err := json.Unmarshal(helper.GetResponseApi(requestUri, timeout), &bid)

	if err != nil {
		log.Fatal(err)
	}

	stringToSave := []byte("Dólar: {" + bid.Value + "}")

	errF := os.WriteFile("cotacao.txt", stringToSave, 0644)

	if errF != nil {
		log.Fatal(errF)
	}

	fmt.Printf("\nValor %v\n", "Dólar: {"+bid.Value+"}")
}
