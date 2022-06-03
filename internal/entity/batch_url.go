package entity

type BatchURLItem struct {
	CorrelationID string
	URL           *URL
}

func NewBatchURLItem(correlationID string, originalURL string, serverAddress string) *BatchURLItem {
	url := NewURL(originalURL, serverAddress)
	return &BatchURLItem{
		CorrelationID: correlationID,
		URL:           url,
	}
}
