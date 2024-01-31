package activequestiondb

type ActiveQuestion struct {
	ID     int `json:"id"`
	UserID int `json:"userId"`
	WordID int `json:"wordId"`
	Type   int `json:"type"`
}
