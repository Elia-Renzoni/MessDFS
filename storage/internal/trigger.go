package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type AuthServiceTrigger struct {
	result QueryResult
}

type QueryResult struct {
	Ack bool `json:"result"`
}

func NewAuthServiceTrigger() *AuthServiceTrigger {
	return &AuthServiceTrigger{}
}

// check the friendship between the owner of the transaction and the owner of the directory
// involved
func (a *AuthServiceTrigger) CheckTransactionOwner(txnUser string, directoryInvolved string) bool {
	if txnUser == directoryInvolved {
		return true
	}

	var partialUrl string = "http://127.0.0.1:8082/friendship?"
	var urlQuery url.Values
	urlQuery.Add("txn", txnUser)
	urlQuery.Add("dir", directoryInvolved)
	var completeUrl string = partialUrl + urlQuery.Encode()

	res, err := http.Get(completeUrl)
	if err != nil {
		return false
	}

	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &a.result)

	if a.result.Ack {
		return true
	}

	return false
}
