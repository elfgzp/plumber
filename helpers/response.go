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

func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func UnauthorizedResponse() JSONResponse {
	return JSONResponse{Code: http.StatusUnauthorized, Msg: "Not Authorized"}
}

