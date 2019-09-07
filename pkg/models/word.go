package models

type Word struct {
	ID      int64
	Profile *Profile
	Token   *Token
	Note    string
}
