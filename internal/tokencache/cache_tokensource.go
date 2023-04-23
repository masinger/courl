package tokencache

import (
	"context"
	"golang.org/x/oauth2"
)

type TokenSaver func(token *oauth2.Token) error
type TokenLoader func() (*oauth2.Token, error)

type TokenSourceCreator func(ctx context.Context, t *oauth2.Token) (oauth2.TokenSource, error)

type tokenCacheSource struct {
	ctx                context.Context
	tokenSourceCreator TokenSourceCreator
	saver              TokenSaver
	loader             TokenLoader
}

func (t tokenCacheSource) Token() (*oauth2.Token, error) {
	cachedToken, err := t.loader()
	if err != nil {
		return cachedToken, err
	}

	src, err := t.tokenSourceCreator(t.ctx, cachedToken)
	if err != nil {
		return nil, err
	}

	freshToken, err := src.Token()
	if err != nil {
		return freshToken, err
	}

	return freshToken, t.saver(freshToken)
}

func NewTokenCache(
	ctx context.Context,
	creator TokenSourceCreator,
	saver TokenSaver,
	loader TokenLoader,
) oauth2.TokenSource {
	return &tokenCacheSource{
		ctx:                ctx,
		tokenSourceCreator: creator,
		saver:              saver,
		loader:             loader,
	}
}
