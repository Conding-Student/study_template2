package response

import "fmt"

type ResponseModel struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Error implements error.
func (r ResponseModel) Error() string {
	return fmt.Sprintf("ResponseModel Error: RetCode=%s, Message=%s", r.RetCode, r.Message)
}
