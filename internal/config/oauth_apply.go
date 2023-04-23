package config

import (
	"github.com/masinger/courl/internal/authprovider"
	"github.com/masinger/courl/internal/util"
)

func (a OAuth) apply(config *authprovider.Config) error {
	util.CopyIfPresent(&a.TokenUrl, &config.TokenUrl)
	util.CopyIfPresent(&a.ClientId, &config.ClientID)
	util.CopyIfPresent(&a.ClientSecret, &config.ClientSecret)

	util.CopyIfPresent(&a.DeviceAuthUrl, &config.DeviceAuthUrl)

	return nil
}
