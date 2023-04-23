package config

import (
	"github.com/masinger/courl/internal/authprovider"
	"net/url"
	"strings"
)

func (c Config) Apply(config *authprovider.Config, urlString string) (*Host, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	for _, hostConfig := range c.Hosts {
		if strings.EqualFold(hostConfig.Host, parsedUrl.Host) {
			appliedHostConfig := hostConfig
			return &appliedHostConfig, hostConfig.apply(config)
		}
	}

	return nil, nil
}
