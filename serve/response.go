package serve

type Response struct {
	Status int                    `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

func NewResponse(status int, data map[string]interface{}) *Response {
	return &Response{
		Status: status,
		Data:   data,
	}
}
