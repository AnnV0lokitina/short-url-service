package entity

// BatchURLItem contains info about URL batch.
type BatchURLItem struct {
	CorrelationID string // // CorrelationID is used to identify request witch send url
	URL           *URL   // URL is used to store information about url
}

// NewBatchURLItem creates new BatchURLItem.
func NewBatchURLItem(correlationID string, originalURL string, serverAddress string) *BatchURLItem {
	url := NewURL(originalURL, serverAddress)
	return &BatchURLItem{
		CorrelationID: correlationID,
		URL:           url,
	}
}
