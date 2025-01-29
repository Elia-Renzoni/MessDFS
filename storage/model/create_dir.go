package model

type CreateDirPayload struct {
	TransactionUser string `json:"txn_user"`
	DirToCreate string `json:"dir_to_create"`
}