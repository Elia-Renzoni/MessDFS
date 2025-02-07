package model


import (
	"net/url"
)

type ReadPayload struct {
	TransactionUser string
	Friend string
	User string
	FileName string
	Query url.Values
}