package handler

// файл содержит описание схем (формата) ввода/вывода
// данные о структурах уровня handler не могут быть переданы в service/entity
import (
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

// JSONRequest json схема для получения 1 url
type JSONRequest struct {
	URL string `json:"url"`
}

// JSONResponse json схема для возвращения 1 url
type JSONResponse struct {
	Result string `json:"result"`
}

// JSONItemRequest json схема для получения списка url
type JSONItemRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// JSONItemResponse json схема для возвращения списка url
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

// JSONListToURLList конвертация полученых данных объекты приложения
func JSONListToURLList(itemInputList []JSONItemRequest, serverAddress string) []*entity.BatchURLItem {
	list := make([]*entity.BatchURLItem, 0, len(itemInputList)) // , len(itemInputList)
	for _, item := range itemInputList {
		urlItem := item.toBatchURLItem(serverAddress)
		list = append(list, urlItem)
	}
	return list
}

// URLListTOJSONList объектов приложения в формат вывода
func URLListTOJSONList(list []*entity.BatchURLItem) []JSONItemResponse {
	outputList := make([]JSONItemResponse, 0, len(list)) // , len(list)
	for _, item := range list {
		i := &JSONItemResponse{}
		i.fromBatchURLItem(item)
		outputList = append(outputList, *i)
	}
	return outputList
}
