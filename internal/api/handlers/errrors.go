package handlers

type ResponseErr struct {
	Error     bool        `json:"error"`
	ErrorText string      `json:"errorText"`
	Data      interface{} `json:"data,omitempty"`
}

func (h *FeedHandlers) ok(data interface{}) *ResponseErr {
	return &ResponseErr{
		Error:     false,
		ErrorText: "",
		Data:      data,
	}
}

func (h *FeedHandlers) error(err error, data interface{}) *ResponseErr {
	return &ResponseErr{
		Error:     true,
		ErrorText: err.Error(),
		Data:      data,
	}
}
