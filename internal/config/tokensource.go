package config

import (
	"context"
	"github.com/masinger/courl/internal/tokencache"
	"golang.org/x/oauth2"
)

func WrapTokenSource(host *Host, ctx context.Context, creator tokencache.TokenSourceCreator) (oauth2.TokenSource, error) {
	if host == nil {
		return creator(ctx, nil)
	}
	return tokencache.WrapTokenSource(host.Host, ctx, creator)
}
