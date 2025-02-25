package abra

import (
	"net/http"
)

type Connector struct {
	httpClient *http.Client
}

func NewAbraConnector(username, password string) *Connector {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	client.Transport = &basicAuthTransport{
		username: username,
		password: password,
		base:     client.Transport,
	}

	return &Connector{httpClient: client}
}

type basicAuthTransport struct {
	username string
	password string
	base     http.RoundTripper
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.username, t.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return t.base.RoundTrip(req)
}
