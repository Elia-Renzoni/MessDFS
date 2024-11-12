package model

import (
	"net/url"
)

type DeletePayload struct {
	user string
	fileName string 
	url url.Values
}