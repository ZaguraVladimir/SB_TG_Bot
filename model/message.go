package model

import (
	"errors"
	"strconv"
	"strings"
)

type Message struct {
	Action     string
	CrCurrency CrCurrency
	Sum        float64
	Text       string
}

func NewMessage(message string) (Message, error) {

	m := Message{
		Action:     "",
		CrCurrency: "",
		Sum:        0,
		Text:       message,
	}
	arr := strings.Split(strings.ToUpper(message), " ")
	count := len(arr)

	m.Action = arr[0]
	if m.Action != "ADD" &&
		m.Action != "SUB" &&
		m.Action != "DEL" &&
		m.Action != "SHOW" {
		return m, errors.New("Неизвестная комманда")
	}

	if m.Action == "ADD" && count != 3 ||
		(m.Action == "SUB" && count != 3) ||
		(m.Action == "DEL" && count != 2) ||
		(m.Action == "SHOW" && count != 1) {
		return m, errors.New("Неверное количество аргументов")
	}

	if count >= 2 {
		m.CrCurrency = CrCurrency(arr[1])
	}

	if count >= 3 {
		ssum := strings.ReplaceAll(arr[2], ",", ".")
		sum, err := strconv.ParseFloat(ssum, 64)
		if err != nil {
			return m, err
		}
		m.Sum = sum
	}

	return m, nil
}
