package clientconfig

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

type Config struct {
	FollowRedirects    bool
	InsecureSkipVerify bool
	IgnoreMissingAuth  bool
}

func (c Config) Create(tokenSource oauth2.TokenSource) (*http.Client, error) {
	baseTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify,
		},
	}

	var client *http.Client

	if tokenSource == nil {
		if !c.IgnoreMissingAuth {
			return nil, fmt.Errorf("no credentials")
		}
		client = &http.Client{
			Transport: baseTransport,
		}
	} else {
		oauthTransport := &oauth2.Transport{
			Source: tokenSource,
			Base:   baseTransport,
		}
		client = &http.Client{
			Transport: oauthTransport,
		}
	}

	if !c.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return client, nil
}

func (c Config) Apply(client *http.Client) error {
	if !c.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	if c.InsecureSkipVerify {
		if client.Transport == nil {
			client.Transport = &http.Transport{}
		}

	}

	return nil
}
