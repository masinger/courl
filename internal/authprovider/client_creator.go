package authprovider

import (
	"context"
	"golang.org/x/oauth2"
)

type tokenSourceProvider func(
	c Config,
	ctx context.Context,
	token *oauth2.Token,
) (oauth2.TokenSource, error)

var tokenSourceProviders = []tokenSourceProvider{}
