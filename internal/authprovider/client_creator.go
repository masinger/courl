package authprovider

import (
	"context"
	"golang.org/x/oauth2"
)

type tokenSourceProvider func(
	c Config,
	ctx context.Context,
) oauth2.TokenSource

var tokenSourceProviders = []tokenSourceProvider{}
