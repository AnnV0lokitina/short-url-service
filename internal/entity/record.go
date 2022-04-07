package entity

type Record struct {
	UserID      uint32 `json:"user_id"`
	Deleted     bool   `json:"deleted"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
