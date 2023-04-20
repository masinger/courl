package config

import "github.com/masinger/courl/internal/authprovider"

func (h Host) apply(config *authprovider.Config) error {
	if h.OAuth != nil {
		return h.OAuth.apply(config)
	}
	return nil
}
