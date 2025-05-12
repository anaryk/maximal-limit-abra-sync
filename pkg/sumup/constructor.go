package sumup

import (
	"net/http"
)

type Connector struct {
	httpClient *http.Client
}

func NewSumUpAPI(token string) *Connector {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	client.Transport = &tokenAuthTransport{
		token: token,
		base:  client.Transport,
	}

	return &Connector{httpClient: client}
}

type tokenAuthTransport struct {
	token string
	base  http.RoundTripper
}

func (t *tokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return t.base.RoundTrip(req)
}
