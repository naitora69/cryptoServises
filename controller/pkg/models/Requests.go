package models

type NewData struct {
	TableName string  `json:"table_name"`
	IDs       []int64 `json:"ids"`
}
