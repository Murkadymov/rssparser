package responses

type ResponseError struct {
	Error     bool        `json:"error"`
	ErrorText string      `json:"errorText"`
	Data      interface{} `json:"data,omitempty"`
}

func OK(data interface{}) *ResponseError {
	return &ResponseError{
		Error:     false,
		ErrorText: "",
		Data:      data,
	}
}

func Error(err error, data interface{}) *ResponseError {
	return &ResponseError{
		Error:     true,
		ErrorText: err.Error(),
		Data:      data,
	}
}
