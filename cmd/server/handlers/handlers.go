package handlers

import (
	"aprokhorov-praktikum/cmd/server/storage"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type HandlerUpdate struct {
	MetricType string
	Storage    storage.Storage
}

func (h HandlerUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		error := fmt.Sprintf("%s, %d, %s", path, len(path), r.URL.Path)
		w.Write([]byte(error))
		return
	}
	switch h.MetricType {
	case "counter":
		current, err := h.Storage.Get("counter", path[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		currentValue := current.(storage.Counter)
		newValue, err := strconv.Atoi(path[4])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, Expected int, got %s", path[4])))
			return
		}
		value := currentValue + storage.Counter(newValue)
		h.Storage.Post(path[3], storage.Counter(value))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Writed %s, previous_value=%d, new value=%d", path[3], current, value)))
	case "gauge":
		newValue, err := strconv.ParseFloat(path[4], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, Expected float, got %s", path[4])))
			return
		}
		value := storage.Gauge(newValue)
		h.Storage.Post(path[3], value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Writed %s, new value=%f", path[3], value)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request, only [gauge, counter] are supported"))
		return
	}
}
