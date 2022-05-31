package entity

// Record Information about URL to save in memory.
type Record struct {
	UserID      uint32 `json:"user_id"`      // user, who save url
	Deleted     bool   `json:"deleted"`      // is url deleted (hidden)
	ShortURL    string `json:"short_url"`    // short url
	OriginalURL string `json:"original_url"` // original url
}
