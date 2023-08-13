package pk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/ratelimit"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	singleRL   ratelimit.Limiter
	bulkRL     ratelimit.Limiter
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	client := &Client{
		apiKey:     apiKey,
		singleRL:   ratelimit.New(1000, ratelimit.Per(60*time.Second)),
		bulkRL:     ratelimit.New(100, ratelimit.Per(60*time.Second)),
		httpClient: nil,
	}
	for _, opt := range opts {
		opt(client)
	}
	if client.httpClient == nil {
		client.httpClient = http.DefaultClient
	}
	return client, nil
}

type Response struct {
	Placekey string `json:"placekey"`
	QueryID  string `json:"query_id"` // filed default value is `"0"`
	Message  string `json:"message"`  // message for error response
}

type getPlacekeyFromGeoRequestParam struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	QueryID   string  `json:"query_id"`
}

type getPlacekeyFromGeoRequest struct {
	Query getPlacekeyFromGeoRequestParam `json:"query"`
}

func (c *Client) req(ctx context.Context, method string, url string, body io.Reader) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.placekey.io/v1/placekey/", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apiKey", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d, body: %s", httpResp.StatusCode, string(respBody))
	}
	resp := &Response{}
	err = json.Unmarshal(respBody, resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// GetPlacekeyFromGeo almost equal like GeoToPlacekey, no `what` part.
func (c *Client) GetPlacekeyFromGeo(
	ctx context.Context, latitude float64, longitude float64,
	queryID string) (*Response, error) {
	_ = c.singleRL.Take()
	apiBody, err := json.Marshal(&getPlacekeyFromGeoRequest{
		Query: getPlacekeyFromGeoRequestParam{
			Latitude:  latitude,
			Longitude: longitude,
			QueryID:   queryID,
		},
	})
	if err != nil {
		return nil, err
	}
	return c.req(ctx, http.MethodPost, "https://api.placekey.io/v1/placekey/", bytes.NewBuffer(apiBody))
}

type getPlacekeyFromAddressRequestParam struct {
	StreetAddress  string `json:"street_address"`
	City           string `json:"city"`
	Region         string `json:"region"`
	PostalCode     string `json:"postal_code"`
	IsoCountryCode string `json:"iso_country_code"`
	QueryID        string `json:"query_id"`
}

type getPlacekeyFromAddressRequest struct {
	Query getPlacekeyFromAddressRequestParam `json:"query"`
}

// GetPlacekeyFromAddress could get both what&where part of placekey.
func (c *Client) GetPlacekeyFromAddress(
	ctx context.Context,
	streetAddress string,
	city string, region string,
	postalCode string, isoCountryCode string,
	queryID string,
) (*Response, error) {
	_ = c.singleRL.Take()
	apiBody, err := json.Marshal(&getPlacekeyFromAddressRequest{
		Query: getPlacekeyFromAddressRequestParam{
			StreetAddress:  streetAddress,
			City:           city,
			Region:         region,
			PostalCode:     postalCode,
			IsoCountryCode: isoCountryCode,
			QueryID:        queryID,
		},
	})
	if err != nil {
		return nil, err
	}
	return c.req(ctx, http.MethodPost, "https://api.placekey.io/v1/placekey/", bytes.NewBuffer(apiBody))
}
