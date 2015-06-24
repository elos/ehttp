package serve

// Error represents a http serving error,
// the status is the HTTP status code, the code
// is a domain specific code. The message is the user
// facing message, and the DevMessage is a developer friendly
// message to help with debugging
type Error struct {
	Status     int    `json:"status"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	DevMessage string `json:"developer_message"`
}

// NewError allocates and returns pointer to a new Error object
func NewError(status, code int, message, devmessage string) *Error {
	return &Error{
		Status:     status,
		Code:       code,
		Message:    message,
		DevMessage: devmessage,
	}
}

// satisfy the error interface
func (e *Error) Error() string {
	return e.DevMessage
}
