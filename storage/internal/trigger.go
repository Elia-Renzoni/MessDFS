package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type AuthServiceTrigger struct {
}


func NewAuthServiceTrigger() *AuthServiceTrigger {
	return &AuthServiceTrigger{}
}

// check the ownership between the owner of the transaction and the owner of the directory
// involved
func (a *AuthServiceTrigger) CheckTransactionOwner(txnUser string, directoryInvolved string) bool {
	var partialUrl string = "http://127.0.0.1:8083/ownership?"
	var urlQuery url.Values
	urlQuery.Add("txn", txnUser)
	urlQuery.Add("dir", directoryInvolved)
	var completeUrl string = partialUrl + urlQuery.Encode()

	res, err := http.Get(completeUrl)
	if err != nil {
		return false
	}

	if res.StatusCode == http.StatusOK {
		return true
	}

	return false
}

func (a *AuthServiceTrigger) CheckFriendship(txnUser string, friend string) bool {
	var partialAuthURL string = "http://127.0.0.1:8083/friendship?"
	var urlQueryFriendship url.Values
	urlQueryFriendship.Add("txn", txnUser)
	urlQueryFriendship.Add("friend", friend)
	var completeUrlFriendship string = partialAuthURL + urlQueryFriendship.Encode()

	res, err := http.Get(completeUrlFriendship)
	if err != nil {
		return false
	}

	if res.StatusCode == http.StatusOK {
		return true
	}

	return false
}