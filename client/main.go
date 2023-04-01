package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	CreateDate string `json:"create_date"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalln("Request Error : ", res.Status)
		return
	}
	payload, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var cotacao Cotacao
	json.Unmarshal(payload, &cotacao)
	err = os.WriteFile("./cotacao.txt", []byte("Dólar: R$"+cotacao.Bid), 0755)

	if err != nil {
		panic(err)
	}
	fmt.Println("Cotação salva com sucesso !")
}
