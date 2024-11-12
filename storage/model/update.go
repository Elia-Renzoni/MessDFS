package model

type UpdatePayload struct {
	QueryType string `json:"query_type"`
	User string `json:"user_name"`
	FileName string `json:"file_name"`
	QueryContent map[string][]string `json:"query_content"`
}