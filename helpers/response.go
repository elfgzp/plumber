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

func Response200(w http.ResponseWriter, msg string, data interface{}) {
	ResponseWithJSON(w, http.StatusOK, OKResponse(msg, data))
}

func Response201(w http.ResponseWriter, msg string, data interface{}) {
	ResponseWithJSON(w, http.StatusCreated, CreatedResponse(msg, data))
}

func Response204(w http.ResponseWriter) {
	ResponseWithJSON(w, http.StatusNoContent, NoContentResponse())
}

func Response500(w http.ResponseWriter) {
	ResponseWithJSON(w, http.StatusInternalServerError, InternalServerErrorResponse())
}

func Response400(w http.ResponseWriter, msg string, data interface{}) {
	ResponseWithJSON(w, http.StatusBadRequest, BadRequestResponse(msg, data))
}

func Response401(w http.ResponseWriter) {
	ResponseWithJSON(w, http.StatusUnauthorized, UnauthorizedResponse())
}

func Response403(w http.ResponseWriter) {
	ResponseWithJSON(w, http.StatusForbidden, PermissionDeniedResponse())
}

func Response404(w http.ResponseWriter, msg string) {
	ResponseWithJSON(w, http.StatusNotFound, NotFoundResponse(msg))
}

func OKResponse(msg string, data interface{}) JSONResponse {
	return JSONResponse{Code: http.StatusOK, Msg: msg, Data: data}
}

func CreatedResponse(msg string, data interface{}) JSONResponse {
	return JSONResponse{Code: http.StatusCreated, Msg: msg, Data: data}
}

func NoContentResponse() JSONResponse {
	return JSONResponse{Code: http.StatusNoContent, Msg: "", Data: nil}
}

func UnauthorizedResponse() JSONResponse {
	return JSONResponse{Code: http.StatusUnauthorized, Msg: "Not Authorized"}
}

func BadRequestResponse(msg string, data interface{}) JSONResponse {
	return JSONResponse{Code: http.StatusBadRequest, Msg: msg, Data: data}
}

func PermissionDeniedResponse() JSONResponse {
	return JSONResponse{Code: http.StatusForbidden, Msg: "Permission denied."}
}

func NotFoundResponse(msg string) JSONResponse {
	if msg == "" {
		msg = "Not Found."
	}
	return JSONResponse{Code: http.StatusNotFound, Msg: msg}
}

func InternalServerErrorResponse() JSONResponse {
	return JSONResponse{Code: http.StatusInternalServerError, Msg: "Internal Server Error."}
}
