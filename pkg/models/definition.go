package models

type Definition struct {
	ID             int64
	Word           *Word
	Classification *Classification
	Definition     string
}
