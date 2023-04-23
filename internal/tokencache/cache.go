package tokencache

import (
	"context"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type TokenCache struct {
	Tokens map[string]*oauth2.Token `yaml:"tokens"`
	path   string
}

func (t TokenCache) Save() (err error) {
	f, err := os.Create(t.path)
	defer func() {
		tempErr := f.Close()
		if tempErr != nil {
			err = tempErr
		}
	}()

	dec := yaml.NewEncoder(f)
	return dec.Encode(t)
}

func LoadCache(path string) (res *TokenCache, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		tmpErr := f.Close()
		if tmpErr != nil {
			err = tmpErr
		}
	}()

	dec := yaml.NewDecoder(f)
	var cache TokenCache
	err = dec.Decode(&cache)
	if err != nil {
		return nil, err
	}
	cache.path = path
	return &cache, nil
}

func LoadDefaultCache() (*TokenCache, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cacheLocation := filepath.Join(userHome, ".courl", ".tokens")
	cache, err := LoadCache(cacheLocation)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err != nil && os.IsNotExist(err) {
		return &TokenCache{
			Tokens: map[string]*oauth2.Token{},
			path:   cacheLocation,
		}, nil
	}
	return cache, err
}

func WrapTokenSource(key string, ctx context.Context, creator TokenSourceCreator) (oauth2.TokenSource, error) {
	cache, err := LoadDefaultCache()
	if err != nil {
		return nil, err
	}

	return NewTokenCache(
		ctx,
		creator,
		func(token *oauth2.Token) error {
			cache.Tokens[key] = token
			return cache.Save()
		},
		func() (*oauth2.Token, error) {
			token, hasKey := cache.Tokens[key]
			if !hasKey {
				return nil, nil
			}
			return token, nil
		},
	), nil
}
