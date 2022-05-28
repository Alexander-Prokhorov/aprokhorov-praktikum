package sender

import (
	"net/http"
	"net/url"
)

type Sender struct {
	Server string
	Port   string
	Url    url.URL
	Client http.Client
}

func (s *Sender) Init() {
	s.Client = http.Client{}
	s.Url = *new(url.URL)
	s.Url.Scheme = "http"
	s.Url.Host = s.Server + ":" + s.Port
}

func (s *Sender) SendMetric(name string, mtype string, value string) error {
	s.Url.Path = "update/" + mtype + "/" + name + "/" + value
	request, err := http.NewRequest(http.MethodPost, s.Url.String(), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	res, err := s.Client.Do(request)
	res.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
