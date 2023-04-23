package authprovider

import (
	"context"
	"golang.org/x/oauth2"
)

type Config struct {
	ClientID      string   `yaml:"clientID"`
	ClientSecret  string   `yaml:"clientSecret"`
	TokenUrl      string   `yaml:"tokenUrl"`
	DeviceAuthUrl string   `yaml:"deviceAuthUrl"`
	Scopes        []string `yaml:"scopes,omitempty"`
}

func (c Config) TokenSource(ctx context.Context, t *oauth2.Token) (oauth2.TokenSource, error) {
	var tokenSource oauth2.TokenSource
	var err error
	for _, provider := range tokenSourceProviders {
		tokenSource, err = provider(c, ctx, t)
		if err != nil {
			return nil, err
		}
		if tokenSource != nil {
			return tokenSource, nil
		}
	}
	return nil, nil
}
