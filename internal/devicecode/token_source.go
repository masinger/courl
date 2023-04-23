package devicecode

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type tokenSource struct {
	config             *Config
	deviceCodeCallback DeviceCodeCallback
	client             *http.Client
	src                oauth2.TokenSource
}

func (c *Config) DeviceCodeTokenSource(
	ctx context.Context,
	httpClient *http.Client,
	callback DeviceCodeCallback,
	token *oauth2.Token,
) oauth2.TokenSource {
	ts := &tokenSource{
		config:             c,
		deviceCodeCallback: callback,
		client:             httpClient,
	}
	if token != nil {
		ts.src = c.TokenSource(ctx, token)
	}
	return oauth2.ReuseTokenSource(token, ts)
}

func (t tokenSource) Token() (*oauth2.Token, error) {
	if t.src != nil {
		wrapped, err := t.src.Token()
		if err == nil {
			return wrapped, nil
		}
	}
	deviceCode, err := t.requestDeviceCode()
	if err != nil {
		return nil, err
	}
	err = t.deviceCodeCallback(deviceCode)
	if err != nil {
		return nil, err
	}
	return t.waitForDeviceAuthorization(deviceCode)
}

func (t tokenSource) completeDeviceAuthorization(deviceCode *DeviceCode) (res *oauth2.Token, err error) {
	params := url.Values{
		"client_secret": {t.config.ClientSecret},
		"client_id":     {t.config.ClientID},
		"grant_type":    {string(t.config.GrantType)},
	}

	params["device_code"] = []string{deviceCode.DeviceCode}

	resp, err := t.client.PostForm(t.config.Endpoint.TokenURL, params)
	if err != nil {
		return nil, err
	}
	defer func() {
		tempErr := resp.Body.Close()
		if tempErr != nil && err == nil {
			err = tempErr
		}
	}()

	// Unmarshal response, checking for errors
	var token tokenOrError
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&token); err != nil {
		return nil, err
	}

	switch token.Error {
	case "":

		// Convert the token into an "oauth2" library token,
		// which doesn't use ExpiresIn.
		return &oauth2.Token{
			AccessToken:  token.AccessToken,
			TokenType:    token.TokenType,
			RefreshToken: token.RefreshToken,
			Expiry:       token.expiry(),
		}, nil
	case "authorization_pending":

	case "slow_down":

		deviceCode.Interval *= 2
	case "access_denied":

		return nil, ErrAccessDenied
	default:

		return nil, fmt.Errorf("authorization failed: %v", token.Error)
	}

	time.Sleep(time.Duration(deviceCode.Interval) * time.Second)
	return nil, nil
}

func (t tokenSource) waitForDeviceAuthorization(deviceCode *DeviceCode) (*oauth2.Token, error) {
	for {
		res, err := t.completeDeviceAuthorization(deviceCode)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
}

func (t tokenSource) requestDeviceCode() (res *DeviceCode, err error) {
	scopes := strings.Join(t.config.Scopes, " ")
	resp, err := t.client.PostForm(t.config.DeviceEndpoint.CodeURL,
		url.Values{"client_id": {t.config.ClientID}, "scope": {scopes}})
	if err != nil {
		return nil, err
	}
	defer func() {
		tempErr := resp.Body.Close()
		if tempErr != nil && err == nil {
			err = tempErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request for device code authorisation returned status %v (%v)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Unmarshal response
	var dcr DeviceCode

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&dcr); err != nil {
		return nil, err
	}

	return &dcr, nil
}
