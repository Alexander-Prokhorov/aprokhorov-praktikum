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
	//body := []byte{'0'}
	var err error
	request, errreq := http.NewRequest(http.MethodPost, s.Url.String(), nil)
	if errreq != nil {
		//fmt.Println("err_req", errreq)
		err = errreq
	}
	request.Header.Set("Content-Type", "text/plain")
	_, errres := s.Client.Do(request)
	if errres != nil {
		//fmt.Println("err_res", errres)
		err = errres
	}
	return err
}
