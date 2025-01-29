package model


import (
	"net/url"
)

type ReadPayload struct {
	TransactionUser string
	User string
	FileName string
	Query url.Values
}