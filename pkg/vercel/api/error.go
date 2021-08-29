package api

// https://vercel.com/docs/api#api-basics/errors

const (
	ErrCodeNotFound = "not_found"
	Forbidden       = "forbidden"
)

type VercelErrorResponse struct {
	Error VercelError `json:"error"`
}

type VercelError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *VercelError) Error() string {
	return e.Message
}

func (e *VercelError) Is(code string) bool {
	return e.Code == code
}
