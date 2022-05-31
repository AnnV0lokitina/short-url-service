package entity

// BatchURLItem Information about URL batch.
type BatchURLItem struct {
	CorrelationID string // correlation id
	URL           *URL   // url list
}

// NewBatchURLItem Create new BatchURLItem.
func NewBatchURLItem(correlationID string, originalURL string, serverAddress string) *BatchURLItem {
	url := NewURL(originalURL, serverAddress)
	return &BatchURLItem{
		CorrelationID: correlationID,
		URL:           url,
	}
}
