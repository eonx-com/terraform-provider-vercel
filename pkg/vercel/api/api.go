package api

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

type Api struct {
	httpClient  *http.Client
	rateLimiter *rate.Limiter

	url       string
	userAgent string
	token     string
}

func New(token string) *Api {
	rateLimiter := rate.NewLimiter(rate.Every(time.Minute/75), 1)

	return &Api{
		httpClient:  &http.Client{},
		rateLimiter: rateLimiter,

		token:     token,
		url:       "https://api.vercel.com",
		userAgent: "eonx-com/terraform-provider-vercel",
	}
}

func (a *Api) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", a.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
}

func (a *Api) getRequestPayload(body interface{}) (io.Reader, error) {
	var payload io.Reader = nil

	if body != nil {
		b, err := json.Marshal(body)

		if err != nil {
			return nil, err
		}

		payload = bytes.NewBuffer(b)
	}

	return payload, nil
}

func (a *Api) Request(method, path string, body interface{}, result interface{}) (*http.Response, *VercelError) {
	payload, err := a.getRequestPayload(body)

	if err != nil {
		return nil, &VercelError{
			Message: "Failed to convert payload request",
		}
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", a.url, path), payload)

	if err != nil {
		return nil, &VercelError{
			Message: err.Error(),
		}
	}

	ctx := context.Background()

	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, &VercelError{
			Message: err.Error(),
		}
	}

	a.setHeaders(req)

	res, err := a.httpClient.Do(req)

	if err != nil {
		return res, &VercelError{
			Message: fmt.Sprint("Unable to perform request: %w", err),
		}
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var verr VercelErrorResponse

		_ = json.NewDecoder(res.Body).Decode(&verr)

		return res, &verr.Error
	}

	if result != nil {
		_ = json.NewDecoder(res.Body).Decode(result)
	}

	return res, nil
}
