package config

import (
	"github.com/masinger/courl/internal/authprovider"
	"net/url"
	"strings"
)

func (c Config) Apply(config *authprovider.Config, urlString string) error {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	for _, hostConfig := range c.Hosts {
		if strings.EqualFold(hostConfig.Host, parsedUrl.Host) {
			return hostConfig.apply(config)
		}
	}

	return nil
}
