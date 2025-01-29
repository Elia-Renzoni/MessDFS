package model

import (
	"net/url"
)

type DeletePayload struct {
	TransactionUser string
	User string
	FileName string 
	Query url.Values
}