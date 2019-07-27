package synology

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	SynoApiAuth = "SYNO.API.Auth"
	SynoApiInfo = "SYNO.API.Info"
)

var supportedAPI = map[string]bool{
	SynoApiAuth:         true,
	SynoApiInfo:         true,
	SYNOFileStationList: true,
}

type ApiInfo struct {
	Name          string
	MaxVersion    int    `json:"maxVersion"`
	MinVersion    int    `json:"minVersion"`
	Path          string `json:"path"`
	RequestFormat string `json:"format"`
}

type apiInfoResponse struct {
	Data    map[string]ApiInfo `json:"data"`
	Success bool               `json:"success"`
}

type loginResponse struct {
	Data struct {
		SID string `json:"sid"`
	} `json:"Data"`
	Success bool `json:"success"`
}

type Client struct {
	Hostname string
	sid      string
	api      map[string]ApiInfo
	log      *log.Logger
	http     *http.Client
}

func (s *Client) prepareRequest(method string, path string, values map[string]string) (*http.Request, error) {
	if s.sid != "" {
		values["_sid"] = s.sid
	}

	url := fmt.Sprintf("%s/webapi/%s", s.Hostname, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for k, v := range values {
		query.Set(k, v)
	}

	req.URL.RawQuery = query.Encode()
	s.log.Println("URL:", req.URL.String())

	return req, nil
}

func (s *Client) do(method string, path string, values map[string]string, out interface{}) error {
	req, err := s.prepareRequest(method, path, values)
	if err != nil {
		return err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if out != nil {
		err = json.Unmarshal(body, out)
		if err != nil {
			s.log.Printf(string(body))
		}
	} else {
		s.log.Printf(string(body))
	}
	return err
}

func (s *Client) download(method string, path string, values map[string]string, w io.Writer) error {
	req, err := s.prepareRequest(method, path, values)
	if err != nil {
		return err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
	return err
}

func (s *Client) GetSupportedApis() map[string]bool {
	return supportedAPI
}

func (s *Client) GetApiInfo() (map[string]ApiInfo, error) {
	params := map[string]string{
		"api":     SynoApiInfo,
		"version": "1",
		"method":  "query",
		"query":   "all",
	}

	resp := &apiInfoResponse{}

	err := s.do("GET", "query.cgi", params, resp)
	return resp.Data, err
}

func (s *Client) login(user, password string) error {
	api, err := s.GetApiInfo()
	if err != nil {
		return err
	}

	loginPath, ok := api[SynoApiAuth]
	if !ok {
		return fmt.Errorf("%s not supported")
	}

	params := map[string]string{
		"api":     SynoApiAuth,
		"version": fmt.Sprintf("%d", loginPath.MaxVersion),
		"method":  "login",
		"account": user,
		"passwd":  password,
		"session": "CLI",
		"format":  "cookie",
	}

	resp := &loginResponse{}

	err = s.do("GET", loginPath.Path, params, resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("login failure")
	}

	s.api = api
	s.sid = resp.Data.SID
	return nil
}

func (s *Client) GetApi(name string) *ApiInfo {
	api, ok := s.api[name]
	if !ok {
		return nil
	}

	api.Name = name
	return &api
}

func NewClient(host string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return &Client{
		Hostname: host,
		http:     client,
	}
}
