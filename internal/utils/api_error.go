package utils

type ApiError struct {
	Message       string `json:"message"`
	StatusCode    int    `json:"statusCode"`
	Status        string `json:"status"`
	Success       bool   `json:"success"`
	IsOperational bool   `json:"isOperational"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(message string, statusCode int) *ApiError {
	status := "Error"

	if statusCode >= 400 && statusCode < 500 {
		status = "Failed"
	}

	return &ApiError{
		Message:       message,
		StatusCode:    statusCode,
		Status:        status,
		Success:       false,
		IsOperational: true,
	}
}
