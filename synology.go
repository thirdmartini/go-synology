package synology

import (
	"io/ioutil"
	"log"
)

type Synology struct {
	client      *Client
	FileStation FileStationService
}

func (s *Synology) WithLogger(log *log.Logger) *Synology {
	s.client.log = log
	return s
}

func (s *Synology) API() (map[string]ApiInfo, error) {
	return s.client.GetApiInfo()
}

func Login(host string, username, password string) (*Synology, error) {
	c := NewClient(host)
	c.log = &log.Logger{}
	c.log.SetOutput(ioutil.Discard)

	err := c.login(username, password)
	if err != nil {
		return nil, err
	}

	return &Synology{
		client:      c,
		FileStation: &FileStationServiceOp{c},
	}, nil
}
