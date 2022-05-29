package handlers

import (
	"aprokhorov-praktikum/cmd/server/storage"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type HandlerUpdate struct {
	//MetricType string
	Storage storage.Storage
}

func (h HandlerUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stringPath := strings.Trim(r.URL.Path, "/")
	slicePath := strings.Split(stringPath, "/")
	var Path struct {
		operation   string
		metricType  string
		metricName  string
		metricValue string
	}
	switch len(slicePath) {
	case 4:
		Path.metricType = slicePath[1]
		Path.metricName = slicePath[2]
		Path.metricValue = slicePath[3]
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		error := fmt.Sprintf("Bad request: %d, %s", len(slicePath), r.URL.Path)
		w.Write([]byte(error))
		return
	}

	switch Path.metricType {
	case "counter":
		current, err := h.Storage.Get("counter", Path.metricName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		currentValue := current.(storage.Counter)
		newValue, err := strconv.Atoi(Path.metricValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, Expected int, got %s", Path.metricValue)))
			return
		}
		value := currentValue + storage.Counter(newValue)
		h.Storage.Post(Path.metricName, storage.Counter(value))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Writed %s, previous_value=%d, new value=%d", Path.metricName, current, value)))
	case "gauge":
		newValue, err := strconv.ParseFloat(Path.metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, Expected float, got %s", Path.metricValue)))
			return
		}
		value := storage.Gauge(newValue)
		h.Storage.Post(Path.metricName, value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Writed %s, new value=%f", Path.metricName, value)))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented yet, only [gauge, counter] are supported"))
		return
	}
}
