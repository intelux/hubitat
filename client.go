package hubitat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Client represents a Hubitat Maker API client.
type Client struct {
	// URL is the base URL of the Maker API. Can be either local or remote.
	URL *url.URL

	// AccessToken is the access token that protects the API.
	AccessToken string

	HTTPClient *http.Client
}

const (
	// EnvHubitatHubURL is the environment variable that contains the Hubitat
	// Hub URL.
	EnvHubitatHubURL = `HUBITAT_HUB_URL`
	// EnvHubitatAccessToken is the environment variable that contains the
	// Hubitat access token.
	EnvHubitatAccessToken = `HUBITAT_ACCESS_TOKEN`
)

// NewClientFromEnv instanciates a new client from the environment.
func NewClientFromEnv() (*Client, error) {
	hubURL := os.Getenv(EnvHubitatHubURL)
	accessToken := os.Getenv(EnvHubitatAccessToken)

	if hubURL == "" {
		return nil, fmt.Errorf("no hub URL set: please set `%s`", EnvHubitatHubURL)
	}

	if accessToken == "" {
		return nil, fmt.Errorf("no access token: please set `%s`", EnvHubitatAccessToken)
	}

	return NewClient(hubURL, accessToken)
}

// NewClient instanciates a new client.
func NewClient(hubURL string, accessToken string) (*Client, error) {
	u, err := url.Parse(hubURL)

	if err != nil {
		return nil, fmt.Errorf("parsing client URL: %s", err)
	}

	return &Client{
		URL:         u,
		AccessToken: accessToken,
		HTTPClient:  http.DefaultClient,
	}, nil
}

// GetDevices list all the devices information.
func (c *Client) GetDevices(ctx context.Context) (Devices, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/apps/api/33/devices/all", nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("requesting devices: %s", err)
	}

	var devices Devices

	if err = json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, fmt.Errorf("decoding devices: %s", err)
	}

	return devices, nil
}

func (c *Client) newRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.url(path).String(), body)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (c *Client) url(path string) *url.URL {
	u := c.URL.ResolveReference(&url.URL{Path: path})

	q := u.Query()
	q.Set("access_token", c.AccessToken)
	u.RawQuery = q.Encode()

	return u
}
