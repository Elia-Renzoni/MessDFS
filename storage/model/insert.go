package model

type InsertPayload struct {
	TransactionUser string `json:"txn_user"`
	QueryType string `json:"query_type"`
	User string `json:"user"`
	FileName string `json:"file_name"`
	QueryContent []string `json:"query_content"`
}