package httputil

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ResponseBody struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJsonOK(w http.ResponseWriter, body ResponseBody) error {
	return WriteJSON(w, http.StatusOK, body)
}

func WriteJSON(w http.ResponseWriter, code int, body ResponseBody) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonBytes)
	return err
}

func ResponseError(w http.ResponseWriter, code int, err error) {
	_ = WriteJSON(w, code, ResponseBody{
		Code:    code,
		Message: strings.Title(err.Error()),
	})
}
