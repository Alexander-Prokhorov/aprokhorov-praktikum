package handlers

import (
	"crypto/hmac"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"aprokhorov-praktikum/cmd/server/storage"
	"aprokhorov-praktikum/internal/hasher"

	"github.com/go-chi/chi"
)

func Post(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Metrics

		req.MType = chi.URLParam(r, "metricType")
		req.ID = chi.URLParam(r, "metricName")
		urlValue := chi.URLParam(r, "metricValue")

		switch req.MType {
		case "counter":
			newValue, err := strconv.ParseInt(urlValue, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("400. Can't parse value to int: %s", urlValue), http.StatusBadRequest)
				return
			}
			req.Delta = &newValue
		case "gauge":
			newValue, err := strconv.ParseFloat(urlValue, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("400. Can't parse value to int: %s", urlValue), http.StatusBadRequest)
				return
			}
			req.Value = &newValue
		}

		err := updateHelper(w, s, &req, "")
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "plain/text")
		http.Error(w, "", http.StatusOK)
	}
}

func updateHelper(w http.ResponseWriter, s storage.Storage, m *Metrics, key string) error {
	switch m.MType {
	case "counter":
		// Проверим валидность хеша
		if m.Hash != "" && key != "" {
			hash := hasher.HashHMAC(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), key)
			if !hmac.Equal([]byte(hash), []byte(m.Hash)) {
				http.Error(w, "invalid hash", http.StatusBadRequest)
				return errors.New("invalid hash")
			}
		}

		// Читаем метрику из базы
		value, err := s.Read(m.MType, m.ID)
		if err != nil {
			// Если метрика не найдена, то устанавливаем счетчик в ноль
			value = storage.Counter(0)
		}

		// Парсим новую метрику из запроса в int64
		if m.Delta == nil {
			http.Error(w, "no delta found", http.StatusBadRequest)
			return errors.New("no delta found")
		}
		newValue := *m.Delta

		// Добавляем значение к прошлому и записываем в сторадж
		resultValue := value.(storage.Counter) + storage.Counter(newValue)
		err = s.Write(m.ID, resultValue)
		if err != nil {
			http.Error(w, "500. Internal Server Error", http.StatusInternalServerError)
			return err
		}

	case "gauge":
		if m.Value == nil {
			http.Error(w, "no value found", http.StatusBadRequest)
			return errors.New("no value found")
		}

		newValue := *m.Value

		// Просто записываем в сторадж
		err := s.Write(m.ID, storage.Gauge(newValue))
		if err != nil {
			http.Error(w, "500. Internal Server Error", http.StatusInternalServerError)
			return err
		}

	default:
		// Вернем NotImplemented, если такой тип еще не поддерживается
		http.Error(w, "501. Not Implemented Yet :)", http.StatusNotImplemented)
		return errors.New("not implemented method")
	}
	return nil
}
