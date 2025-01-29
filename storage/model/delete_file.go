package model 

type DeleteFilePayload struct {
	TransactionUser string
	FileToDelete string
	DirName string
}