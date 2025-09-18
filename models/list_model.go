package models

type ListEntry struct {
	ID         int64
	InsertTime string
	Group      string
	Title      string
	Priority   int8
	Status     int8
}
