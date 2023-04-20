package authprovider

import (
	"context"
	"github.com/masinger/courl/internal/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func init() {
	tokenSourceProviders = append(tokenSourceProviders, func(c Config, ctx context.Context) oauth2.TokenSource {
		if !util.AllStringsPresent(&c.TokenURL, &c.ClientID, &c.ClientSecret) {
			return nil
		}

		clientConfig := clientcredentials.Config{
			TokenURL:     c.TokenURL,
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			Scopes:       c.Scopes,
		}

		return clientConfig.TokenSource(ctx)
	})
}
