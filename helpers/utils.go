package helper

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  int         `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
