package sender

import (
	"net/http"
	"net/url"
)

type Sender struct {
	Server string
	Port   string
	URL    url.URL
	Client http.Client
}

func (s *Sender) Init() {
	s.Client = http.Client{}
	s.URL = *new(url.URL)
	s.URL.Scheme = "http"
	s.URL.Host = s.Server + ":" + s.Port
}

func (s *Sender) SendMetric(name string, mtype string, value string) error {
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
