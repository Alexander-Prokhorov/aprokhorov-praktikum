package handlers

import (
	"aprokhorov-praktikum/cmd/server/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func Post(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")
		switch metricType {
		case "counter":
			// Читаем метрику из базы
			value, err := s.Read(metricType, metricName)
			if err != nil {
				// Если метрика не найдена, то устанавливаем счетчик в ноль
				value = storage.Counter(0)
			}

			// Парсим новую метрику из запроса в int64
			newValue, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				// Если не можем распарсить возвращаем ошибку
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400. Bad Request"))
			}

			// Добавляем значение к прошлому и записываем в сторадж
			resultValue := value.(storage.Counter) + storage.Counter(newValue)
			s.Write(metricName, resultValue)

		case "gauge":
			newValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				// Если не можем распарсить возвращаем ошибку
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400. Bad Request"))
			}

			// Просто записываем в сторадж
			s.Write(metricName, storage.Gauge(newValue))

		default:
			// Вернем NotImplemented, если такой тип еще не поддерживается
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("501. Not Implemented Yet :)"))
		}
	}
}
