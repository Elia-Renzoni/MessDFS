package model

type InsertPayload struct {
	QueryType string `json:"query_type"`
	User string `json:"user"`
	FileName string `json:"file_name"`
	QueryContent []string `json:"query_content"`
}