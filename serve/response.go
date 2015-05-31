package serve

type Response struct {
	Status uint64                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

func NewResponse(status uint64, data map[string]interface{}) *Response {
	return &Response{
		Status: status,
		Data:   data,
	}
}
