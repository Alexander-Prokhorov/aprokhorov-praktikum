package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"aprokhorov-praktikum/internal/ccrypto"
	"aprokhorov-praktikum/internal/hasher"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// Agent Sender Data.
type HTTPSender struct {
	Server   string
	Port     string
	URL      url.URL
	Client   http.Client
	sourceIP string
}

// Create and init new Agent Sender.
func NewAgentSender(address string) *HTTPSender {
	const (
		defaultTimeout         = 5
		defaultMaxIdleConns    = 5
		defaultIdleConnTimeout = 5
	)

	var s HTTPSender

	s.Client = http.Client{
		Timeout: time.Second * defaultTimeout,
		Transport: &http.Transport{
			MaxIdleConns:    defaultMaxIdleConns,
			IdleConnTimeout: defaultIdleConnTimeout,
		},
	}
	s.URL = url.URL{}
	s.URL.Scheme = "http"
	s.URL.Host = address

	return &s
}

// Init Source Address from local host addresses.
// Prefer Global->Private->Loopback->LinkLocal.
// Prefer biggest IP, if found several.
func (s *HTTPSender) InitSourceAddress() (string, error) {
	ip, err := getOutputAddr()

	s.sourceIP = ip.String()

	return s.sourceIP, err
}

// Send Metric to Server.
func (s *HTTPSender) SendMetricSingle(
	ctx context.Context,
	mtype string,
	name string,
	value string,
	hashKey string,
	pubKey *ccrypto.PublicKey,
) error {
	return s.SendMetricJSON(mtype, name, value, hashKey, pubKey)
}

// Send Batch of metrics to Server.
func (s *HTTPSender) SendMetricBatch(
	ctx context.Context,
	metrics map[string]map[string]string,
	hashKey string,
	pubKey *ccrypto.PublicKey,
) error {
	return s.SendMetricJSONBatch(metrics, hashKey, pubKey)
}

// Send Metrics by POST-req to Server url-encoded.
func (s *HTTPSender) SendMetricURL(mtype string, name string, value string, key string) error {
	s.URL.Path = "update/" + mtype + "/" + name + "/" + value

	request, err := http.NewRequest(http.MethodPost, s.URL.String(), nil)
	if err != nil {
		return err
	}

	request.Header.Set("X-Real-IP", s.sourceIP)
	request.Header.Set("Content-Type", "text/plain")

	res, err := s.Client.Do(request)
	if err != nil {
		return err
	}

	return res.Body.Close()
}

// Send Metric by POST-req for Server JSON-body.
func (s *HTTPSender) SendMetricJSON(
	mtype string,
	name string,
	value string,
	hashKey string,
	pubKey *ccrypto.PublicKey,
) error {
	s.URL.Path = "update/"

	req, err := s.helperSendMetricJSON(mtype, name, value, hashKey)
	if err != nil {
		return err
	}

	jReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if pubKey != nil {
		jReq, err = pubKey.Encrypt(jReq)
		if err != nil {
			return err
		}
	}

	request, err := http.NewRequest(http.MethodPost, s.URL.String(), bytes.NewBuffer(jReq))
	if err != nil {
		return err
	}

	request.Header.Set("X-Real-IP", s.sourceIP)
	request.Header.Set("Content-Type", "application/json")

	res, err := s.Client.Do(request)
	if err != nil {
		return err
	}

	return res.Body.Close()
}

func (s *HTTPSender) helperSendMetricJSON(mtype string, name string, value string, key string) (Metrics, error) {
	const (
		bitSize = 64
		base    = 10
	)

	var req Metrics
	req.ID = name
	req.MType = mtype

	switch mtype {
	case Counter:
		rValue, err := strconv.ParseInt(value, base, bitSize)
		if err != nil {
			return Metrics{}, err
		}

		req.Delta = &rValue
		if key != "" {
			req.Hash = hasher.HashHMAC(fmt.Sprintf("%s:counter:%d", name, rValue), key)
		}
	case Gauge:
		rValue, err := strconv.ParseFloat(value, bitSize)
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

// Send batch of Metrics by POST-req with JSON-body.
func (s *HTTPSender) SendMetricJSONBatch(
	metrics map[string]map[string]string,
	hashKey string,
	pubKey *ccrypto.PublicKey,
) error {
	s.URL.Path = "updates/"

	batchReq := make([]Metrics, 0)

	for metricType, metric := range metrics {
		for metricName, metricValue := range metric {
			req, err := s.helperSendMetricJSON(metricType, metricName, metricValue, hashKey)
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

	if pubKey != nil {
		jReq, err = pubKey.Encrypt(jReq)
		if err != nil {
			return err
		}
	}

	request, err := http.NewRequest(http.MethodPost, s.URL.String(), bytes.NewBuffer(jReq))
	if err != nil {
		return err
	}

	request.Header.Set("X-Real-IP", s.sourceIP)
	request.Header.Set("Content-Type", "application/json")

	res, err := s.Client.Do(request)
	if err != nil {
		return err
	}

	return res.Body.Close()
}

func getOutputAddr() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var ip net.IP
	var gIP, pIP, lIP, llIP []net.IP
	for _, addr := range addrs {
		ip = addr.(*net.IPNet).IP
		switch {
		case ip.IsGlobalUnicast() && !ip.IsPrivate():
			gIP = append(gIP, ip)
		case ip.IsPrivate():
			pIP = append(pIP, ip)
		case ip.IsLoopback():
			lIP = append(lIP, ip)
		case ip.IsLinkLocalUnicast():
			llIP = append(llIP, ip)
		}
	}
	switch {
	case len(gIP) != 0:
		return biggestIP(gIP), nil
	case len(pIP) != 0:
		return biggestIP(pIP), nil
	case len(lIP) != 0:
		return biggestIP(lIP), nil
	case len(llIP) != 0:
		return biggestIP(llIP), nil
	default:
		return ip, nil
	}
}

func biggestIP(ips []net.IP) net.IP {
	var bIP net.IP
	switch len(ips) {
	case 0:
		return nil
	case 1:
		bIP = ips[0]
	default:
		bIP = ips[0]
		for _, ip := range ips[1:] {
			if bytes.Compare(bIP, ip) < 0 {
				bIP = ip
			}
		}
	}
	return bIP
}
