package util

import (
	"encoding/json"
	"net/http"
)

//ErrorResponse composed by a status code and an error message
type ErrorResponse struct {
	Status int    `json:"status"`
	Err    string `json:"error"`
}

func EncodeError(response http.ResponseWriter, statusCode int, message string) {
	response.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(response).Encode(ErrorResponse{Status: statusCode, Err: message})
}
