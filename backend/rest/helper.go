package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Code  ...
type Code string

// Response response serializer util
type Response struct {
	Code    Code        `json:"code,omitempty"`
	Status  int         `json:"-"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Count   *int64      `json:"count,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ServeJSON serves json to http client
func (r *Response) ServeJSON(w http.ResponseWriter) error {
	resp := &Response{
		Code:    r.Code,
		Status:  r.Status,
		Message: r.Message,
		Data:    r.Data,
		Errors:  r.Errors,
		Count:   r.Count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

// ServeJSON a utility func which serves json to http client
func ServeJSON(w http.ResponseWriter, code Code, status int, message string, data interface{}, count *int64, errors interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var resp interface{}

	resp = &Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
		Errors:  errors,
		Count:   count,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

type reqList struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

func getPager(r *http.Request) *reqList {
	rs := &reqList{}

	l := r.URL.Query().Get("limit")
	if li, err := strconv.Atoi(l); err == nil {
		if li <= 0 || li > 25 {
			li = 25
		}
		rs.Limit = int64(li)
	} else {
		rs.Limit = int64(25)
	}

	p := r.URL.Query().Get("page")
	if pi, err := strconv.Atoi(p); err == nil {
		if pi <= 0 {
			pi = 1
		}
		rs.Page = int64(pi)
	} else {
		rs.Page = int64(1)
	}

	return rs
}
