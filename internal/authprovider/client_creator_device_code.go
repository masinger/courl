package authprovider

import (
	"context"
	"fmt"
	"github.com/masinger/courl/internal/devicecode"
	"github.com/masinger/courl/internal/util"
	"golang.org/x/oauth2"
	"net/http"
)

func init() {
	tokenSourceProviders = append(tokenSourceProviders, func(c Config, ctx context.Context, t *oauth2.Token) (oauth2.TokenSource, error) {
		if !util.AllStringsPresent(&c.TokenUrl, &c.ClientID, &c.DeviceAuthUrl) {
			return nil, nil
		}
		config := devicecode.Config{
			Config: oauth2.Config{
				ClientID:     c.ClientID,
				ClientSecret: c.ClientSecret,
				Endpoint: oauth2.Endpoint{
					TokenURL: c.TokenUrl,
				},
				Scopes: c.Scopes,
			},
			DeviceEndpoint: devicecode.DeviceEndpoint{CodeURL: c.DeviceAuthUrl},
			GrantType:      devicecode.DeviceGrantTypeIetf,
		}

		return config.DeviceCodeTokenSource(
			ctx,
			http.DefaultClient,
			func(deviceCode *devicecode.DeviceCode) error {
				fmt.Printf("Visit: %v\n", deviceCode.VerificationUriComplete)
				return nil
			},
			t,
		), nil
	})
}
