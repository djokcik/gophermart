package handler

import "net/http"

func NotImplementedHandler(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "not implemented", http.StatusNotImplemented)
}
