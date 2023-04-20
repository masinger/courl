package authprovider

import (
	"context"
	"golang.org/x/oauth2"
)

type Config struct {
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	TokenURL     string   `yaml:"tokenURL"`
	Scopes       []string `yaml:"scopes,omitempty"`
}

func (c Config) TokenSource(ctx context.Context) oauth2.TokenSource {
	var tokenSource oauth2.TokenSource
	for _, provider := range tokenSourceProviders {
		tokenSource = provider(c, ctx)
		if tokenSource != nil {
			return tokenSource
		}
	}
	return nil
}
