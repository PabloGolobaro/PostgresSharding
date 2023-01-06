package models

type Item struct {
	ID   int64  `db:"id"`
	Code string `db:"code"`
}
