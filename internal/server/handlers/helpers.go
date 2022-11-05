package handlers

import (
	"context"
	"crypto/hmac"
	"errors"
	"fmt"
	"net/http"

	"aprokhorov-praktikum/internal/hasher"
	"aprokhorov-praktikum/internal/storage"
)

const (
	Counter     = "counter"
	Gauge       = "gauge"
	contentJSON = "application/json"
)

func updateHelper(ctx context.Context, w http.ResponseWriter, s storage.Storage, m *Metrics, key string) error {
	switch m.MType {
	case Counter:
		// Проверим валидность хеша
		if m.Hash != "" && key != "" {
			hash := hasher.HashHMAC(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), key)
			if !hmac.Equal([]byte(hash), []byte(m.Hash)) {
				http.Error(w, "invalid hash", http.StatusBadRequest)

				return errors.New("invalid hash")
			}
		}

		// Читаем метрику из базы
		value, err := s.Read(ctx, m.MType, m.ID)
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
		oldValue, ok := value.(storage.Counter)
		if !ok {
			http.Error(w, "500. Internal Server Error", http.StatusInternalServerError)

			return errors.New("cannot make assertion (storage.Counter)")
		}

		resultValue := oldValue + storage.Counter(newValue)

		err = s.Write(ctx, m.ID, resultValue)
		if err != nil {
			http.Error(w, "500. Internal Server Error", http.StatusInternalServerError)

			return err
		}

	case Gauge:
		if m.Value == nil {
			http.Error(w, "no value found", http.StatusBadRequest)

			return errors.New("no value found")
		}

		newValue := *m.Value

		// Просто записываем в сторадж
		err := s.Write(ctx, m.ID, storage.Gauge(newValue))
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

func readHelper(ctx context.Context, w http.ResponseWriter, s storage.Reader, m *Metrics, key string) error {
	var hashString string

	value, err := s.Read(ctx, m.MType, m.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("404. Not Found. %s", err), http.StatusNotFound)

		return err
	}

	switch data := value.(type) {
	case storage.Counter:
		respond := int64(data)
		m.Delta = &respond
		hashString = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case storage.Gauge:
		respond := float64(data)
		m.Value = &respond
		hashString = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		http.Error(w, fmt.Sprintf("500. Internal Server. %s", err), http.StatusInternalServerError)

		return err
	}

	if key != "" {
		m.Hash = hasher.HashHMAC(hashString, key)
	}

	return nil
}
