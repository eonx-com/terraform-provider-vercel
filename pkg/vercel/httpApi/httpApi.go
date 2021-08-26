package httpApi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type API interface {
	Request(method string, path string, body interface{}) (*http.Response, error)
}

type Api struct {
	httpClient  *http.Client
	rateLimiter *rate.Limiter

	url       string
	userAgent string
	token     string
}

func New(token string) API {
	return &Api{
		httpClient:  &http.Client{},
		rateLimiter: rate.NewLimiter(rate.Every(800*time.Millisecond), 1),

		url:       "https://api.vercel.com",
		userAgent: "eonx-com/terraform-provider-vercel",
		token:     token,
	}
}

// https://vercel.com/docs/api#api-basics/errors
type VercelError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Api) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
}

func (c *Api) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()

	err := c.rateLimiter.Wait(ctx)

	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		defer res.Body.Close()

		var x map[string]interface{}
		_ = json.NewDecoder(res.Body).Decode(&x)

		return res, fmt.Errorf("error during http request: %+v", x)
	}

	return res, nil
}

func (c *Api) Request(method string, path string, body interface{}) (*http.Response, error) {
	var payload io.Reader = nil

	if body != nil {
		b, err := json.Marshal(body)

		if err != nil {
			return nil, err
		}

		payload = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.url, path), payload)

	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)

	if err != nil {
		return res, fmt.Errorf("unable to request resource: [%s] %s with payload {%+v}: %w", method, path, payload, err)
	}

	return res, nil

}
