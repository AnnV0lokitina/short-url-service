package handler

// file contains description of I/O schemes (format)
// data about handler level structures cannot be passed to service/entity
import (
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

// JSONRequest json schema to get 1 url
type JSONRequest struct {
	URL string `json:"url"`
}

// JSONResponse json schema to return 1 url
type JSONResponse struct {
	Result string `json:"result"`
}

// JSONItemRequest json schema to get a list of urls
type JSONItemRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// JSONItemResponse json schema to return a list of urls
type JSONItemResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (ii JSONItemRequest) toBatchURLItem(serverAddress string) *entity.BatchURLItem {
	return entity.NewBatchURLItem(
		ii.CorrelationID,
		ii.OriginalURL,
		serverAddress,
	)
}

func (io *JSONItemResponse) fromBatchURLItem(item *entity.BatchURLItem) {
	io.CorrelationID = item.CorrelationID
	io.ShortURL = item.URL.Short
}

// JSONListToURLList converting received data application objects
func JSONListToURLList(itemInputList []JSONItemRequest, serverAddress string) []*entity.BatchURLItem {
	list := make([]*entity.BatchURLItem, 0, len(itemInputList)) // , len(itemInputList)
	for _, item := range itemInputList {
		urlItem := item.toBatchURLItem(serverAddress)
		list = append(list, urlItem)
	}
	return list
}

// URLListTOJSONList application objects to output format
func URLListTOJSONList(list []*entity.BatchURLItem) []JSONItemResponse {
	outputList := make([]JSONItemResponse, 0, len(list)) // , len(list)
	for _, item := range list {
		i := &JSONItemResponse{}
		i.fromBatchURLItem(item)
		outputList = append(outputList, *i)
	}
	return outputList
}
