package model


import (
	"net/url"
)

type ReadPayload struct {
	User string
	FileName string
	url url.Values
}