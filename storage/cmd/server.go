package main

import (
	"net/http"
	"storage/storage"
)

func main() {
	storageService := storage.NewStorage(":8081")
	storageService.Start()
}