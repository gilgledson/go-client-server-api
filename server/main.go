package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type USDBRL struct {
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

type CotacaoDolar struct {
	Usdbrl USDBRL `json:"USDBRL"`
}

func main() {

	http.HandleFunc("/cotacao", getCotacaoDolar)
	http.ListenAndServe(":8080", nil)

}
func _appError(w http.ResponseWriter, err error) {
	log.Fatalln(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
func getCotacaoDolar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		_appError(w, err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		_appError(w, err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		_appError(w, err)
		return
	}

	cotacao := CotacaoDolar{}
	json.Unmarshal(body, &cotacao)
	err = saveOnDatabase(cotacao)
	if err != nil {
		_appError(w, err)
		return
	}

	response, err := json.Marshal(cotacao.Usdbrl)
	if err != nil {
		_appError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func saveOnDatabase(cotacao CotacaoDolar) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})

	if err != nil {
		return err
	}

	db.AutoMigrate(&USDBRL{})
	db = db.WithContext(queryCtx)
	err = db.Create(&cotacao.Usdbrl).Error
	if err != nil {
		return err
	}
	return nil
}
