package authprovider

import (
	"context"
	"github.com/masinger/courl/internal/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func init() {
	tokenSourceProviders = append(tokenSourceProviders, func(c Config, ctx context.Context, _ *oauth2.Token) (oauth2.TokenSource, error) {
		if !util.AllStringsPresent(&c.TokenUrl, &c.ClientID, &c.ClientSecret) {
			return nil, nil
		}

		clientConfig := clientcredentials.Config{
			TokenURL:     c.TokenUrl,
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			Scopes:       c.Scopes,
		}
		return clientConfig.TokenSource(ctx), nil
	})
}
