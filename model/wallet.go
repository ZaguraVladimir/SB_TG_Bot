package model

import (
	"encoding/json"
	"fmt"
	"log"
	http "net/http"
	"strings"
)

type CrCurrency string
type Wallet map[CrCurrency]float64

type bnResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
	Code   int64   `json:"code,string"`
	Msg    string  `json:"msg"`
}

func (w *Wallet) Processing(message Message) string {

	if _, ok := (*w)[message.CrCurrency]; !ok && (message.Action == "SUB" || message.Action == "DEL") {
		return fmt.Sprintf("Операция %s невозможна, валюта %s отсутствует в кощельке", message.Action, message.CrCurrency)
	}

	if len(*w) == 0 && (message.Action == "SHOW") {
		return "Нет средств на балансе"
	}

	if message.Sum > (*w)[message.CrCurrency] && (message.Action == "SUB") {
		return "На балансе недостаточно средств"
	}

	result := "Выполнено"

	switch message.Action {
	case "ADD":
		(*w)[message.CrCurrency] += message.Sum
	case "SUB":
		(*w)[message.CrCurrency] -= message.Sum
	case "DEL":
		delete(*w, message.CrCurrency)
	case "SHOW":
		builder := strings.Builder{}
		builder.WriteString("Баланс:\n")

		totalUSD := .0
		priceRUB := getPrice(CrCurrency("USDT"), CrCurrency("RUB"))
		for key, sum := range *w {
			priceUSD := getPrice(key, CrCurrency("USDT"))
			sumUSD := sum * priceUSD
			sumRUB := sumUSD * priceRUB
			builder.WriteString(fmt.Sprintf("%s: %f (%.2f$ %.2f₽)\n", key, sum, sumUSD, sumRUB))
			totalUSD += sumUSD
		}
		builder.WriteString(fmt.Sprintf("\nВсего в $: %.2f", totalUSD))
		builder.WriteString(fmt.Sprintf("\nВсего в ₽: %.2f", totalUSD*priceRUB))
		result = builder.String()
	}
	return result
}

func getPrice(currency1 CrCurrency, currency2 CrCurrency) float64 {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", currency1, currency2)
	response, err := http.Get(url)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	var bnJson bnResponse
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&bnJson); err != nil {
		log.Print(err)
	}
	return bnJson.Price
}
