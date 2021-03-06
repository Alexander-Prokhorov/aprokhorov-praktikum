package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"aprokhorov-praktikum/internal/hasher"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type Sender struct {
	Server string
	Port   string
	URL    url.URL
	Client http.Client
}

func NewAgentSender(address string) *Sender {
	var s Sender
	s.Client = http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    5,
			IdleConnTimeout: 5,
		},
	}
	s.URL = *new(url.URL)
	s.URL.Scheme = "http"
	s.URL.Host = address

	return &s
}

func (s *Sender) SendMetric(mtype string, name string, value string, key string) error {
	return s.SendMetricJSON(mtype, name, value, key)
}

func (s *Sender) SendMetricBatch(metrics map[string]map[string]string, key string) error {
	return s.SendMetricJSONBatch(metrics, key)
}

func (s *Sender) SendMetricURL(mtype string, name string, value string, key string) error {
	s.URL.Path = "update/" + mtype + "/" + name + "/" + value
	request, err := http.NewRequest(http.MethodPost, s.URL.String(), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	res, err := s.Client.Do(request)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}

func (s *Sender) SendMetricJSON(mtype string, name string, value string, key string) error {
	s.URL.Path = "update/"

	req, err := s.helperSendMetricJSON(mtype, name, value, key)
	if err != nil {
		return err
	}

	jReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, s.URL.String(), bytes.NewBuffer(jReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	res, err := s.Client.Do(request)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}

func (s *Sender) helperSendMetricJSON(mtype string, name string, value string, key string) (Metrics, error) {
	var req Metrics
	req.ID = name
	req.MType = mtype

	switch mtype {
	case Counter:
		rValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metrics{}, err
		}
		req.Delta = &rValue
		if key != "" {
			req.Hash = hasher.HashHMAC(fmt.Sprintf("%s:counter:%d", name, rValue), key)
		}
	case Gauge:
		rValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metrics{}, err
		}
		req.Value = &rValue
		if key != "" {
			req.Hash = hasher.HashHMAC(fmt.Sprintf("%s:gauge:%f", name, rValue), key)
		}
	}
	return req, nil
}

func (s *Sender) SendMetricJSONBatch(metrics map[string]map[string]string, key string) error {
	s.URL.Path = "updates/"

	batchReq := make([]Metrics, 0)

	for metricType, metric := range metrics {
		for metricName, metricValue := range metric {
			req, err := s.helperSendMetricJSON(metricType, metricName, metricValue, key)
			if err != nil {
				return err
			}
			batchReq = append(batchReq, req)
		}
	}

	jReq, err := json.Marshal(batchReq)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, s.URL.String(), bytes.NewBuffer(jReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	res, err := s.Client.Do(request)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}
