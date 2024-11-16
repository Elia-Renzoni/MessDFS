package model

import (
	"net/url"
)

type DeletePayload struct {
	User string
	FileName string 
	Query url.Values
}