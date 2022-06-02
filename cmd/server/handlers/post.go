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
				_, err := w.Write([]byte("400. Bad Request"))
				if err != nil {
					panic(err)
				}
			}

			// Добавляем значение к прошлому и записываем в сторадж
			resultValue := value.(storage.Counter) + storage.Counter(newValue)
			err = s.Write(metricName, resultValue)
			if err != nil {
				panic(err)
			}

		case "gauge":
			newValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				// Если не можем распарсить возвращаем ошибку
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("400. Bad Request"))
				if err != nil {
					panic(err)
				}
			}

			// Просто записываем в сторадж
			err = s.Write(metricName, storage.Gauge(newValue))
			if err != nil {
				panic(err)
			}

		default:
			// Вернем NotImplemented, если такой тип еще не поддерживается
			w.WriteHeader(http.StatusNotImplemented)
			_, err := w.Write([]byte("501. Not Implemented Yet :)"))
			if err != nil {
				panic(err)
			}
		}
	}
}
