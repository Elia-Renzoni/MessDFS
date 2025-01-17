package main

import (
	ss "storageservice/internal"
)

func main() {
	storageService := ss.NewStorage(":8081")
	storageService.Start()
}
