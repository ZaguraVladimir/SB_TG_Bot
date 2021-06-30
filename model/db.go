package model

type ChatID int64
type DB map[ChatID]Wallet

func (с ChatID) ID() int64 {
	return int64(с)
}
func (db *DB) AddID(id ChatID) {
	if _, ok := (*db)[id]; !ok {
		(*db)[id] = make(Wallet)
	}
}
