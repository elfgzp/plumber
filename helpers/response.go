package helpers

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type PagedData struct {
	Total  int           `json:"total"`
	Page   int           `json:"page"`
	Limit  int           `json:"limit"`
	Result []interface{} `json:"result"`
}

func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func UnauthorizedResponse() JSONResponse {
	return JSONResponse{Code: http.StatusUnauthorized, Msg: "Not Authorized"}
}

func InternalServerErrorResponse() JSONResponse {
	return JSONResponse{Code: http.StatusInternalServerError, Msg: "Internal Server Error."}
}
