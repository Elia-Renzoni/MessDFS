package model

import (
	"net/url"
)

type DeletePayload struct {
	User string
	FileName string 
	Url url.Values
}