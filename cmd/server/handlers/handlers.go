package handlers

import (
	"aprokhorov-praktikum/cmd/server/storage"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type HandlerUpdate struct {
	Metric_type string
	Storage     storage.Storage
}

func (h HandlerUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		error := fmt.Sprintf("%s, %d, %s", path, len(path), r.URL.Path)
		w.Write([]byte(error))
		return
	}
	switch h.Metric_type {
	case "counter":
		current, err := h.Storage.Get("counter", path[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		current_value := current.(storage.Counter)
		new_value, err := strconv.Atoi(path[4])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, Expected int, got %s", path[4])))
			return
		}
		value := current_value + storage.Counter(new_value)
		h.Storage.Post(path[3], storage.Counter(value))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Writed %s, previous_value=%d, new value=%d", path[3], current, value)))
	case "gauge":
		new_value, err := strconv.ParseFloat(path[4], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, Expected float, got %s", path[4])))
			return
		}
		value := storage.Gauge(new_value)
		h.Storage.Post(path[3], value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Writed %s, new value=%f", path[3], value)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request, only [gauge, counter] are supported"))
		return
	}
}
