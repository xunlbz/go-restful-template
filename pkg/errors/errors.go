package errors

import (
	"encoding/json"
)

// ReponseError Common Error
type ReponseError struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (re *ReponseError) Error() string {
	err, _ := json.Marshal(re)
	return string(err)
}

// New create a new ReponseError
func New(code int, msg string) *ReponseError {
	return &ReponseError{Code: code, Msg: msg}
}
