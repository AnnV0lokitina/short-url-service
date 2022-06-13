package entity

// Record Information about URL to save in memory.
type Record struct {
	UserID      uint32 `json:"user_id"` // owner's user ID
	Deleted     bool   `json:"deleted"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
