package main

import (
	"net/http"
	ss "storageservice/internal"
)

func main() {
	storageService := ss.NewStorage(":8081")
	storageService.Start()
}