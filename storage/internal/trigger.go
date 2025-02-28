package internal

import (
	"fmt"
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
// this method is called by every endpoint, except for the read operations.
func (a *AuthServiceTrigger) CheckTransactionOwner(txnUser string, directoryInvolved string) bool {
	var partialUrl string = "http://127.0.0.1:8083/ownership?"
	fmt.Printf("%s", txnUser)
	urlQuery := url.Values{}
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

// this method is responsible for checking the frienship beetwen users
// the method perform a GET call to the auth microservice
func (a *AuthServiceTrigger) CheckFriendship(txnUser string, friend string) bool {
	if friend == txnUser {
		return true
	}

	var partialAuthURL string = "http://127.0.0.1:8083/friendship?"
	fmt.Printf("%s", txnUser)
	urlQueryFriendship := url.Values{}
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

func (a *AuthServiceTrigger) DeleteDirectoryInAuth(dirname string) bool {
	var partialAuthURL string = "http://127.0.0.1:8083/delete-dir?"
	urlQuery := url.Values{}

	urlQuery.Add("directory", dirname)
	var completeDeleteUrl string = partialAuthURL + urlQuery.Encode()

	fmt.Println(completeDeleteUrl)

	request, errReq := http.NewRequest("DELETE", completeDeleteUrl, nil)
	if errReq != nil {
		return false
	}

	client := http.Client{}

	res, err := client.Do(request)
	if err != nil {
		return false
	}

	if res.StatusCode == http.StatusOK {
		return true
	}

	return false
}
