package ginhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

// EncodeRequestFunc encodes the passed request object into the HTTP request
// object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward EncodeRequestFunc could be something that JSON
// encodes the object directly to the request body.
type EncodeRetryableRequestFunc[req any] func(context.Context, *retryablehttp.Request, *req) error

// CreateRequestFunc creates an outgoing HTTP request based on the passed
// request object. It's designed to be used in HTTP clients, for client-side
// endpoints. It's a more powerful version of EncodeRequestFunc, and can be used
// if more fine-grained control of the HTTP request is required.
type CreateRetryableRequestFunc[req any] func(context.Context, *req) (*retryablehttp.Request, error)

// encodeHttpGenericRequest is a transport/http.EncodeRequestFunc that
// SON-encodes any request to the request body. Primarily useful in a client.
func encodeHttpGenericRetryableRequest[req any](_ context.Context, r *retryablehttp.Request, request *req) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	r.Body = io.NopCloser(bytes.NewReader(b))

	return nil
}

type RetryableClient[req, resp any] struct {
	client *retryablehttp.Client
	req    CreateRetryableRequestFunc[req]
	dec    DecodeResponseFunc[resp]
}

// NewClient constructs a usable Client for a single remote method.
func NewRetryableClient[req, resp any](client *retryablehttp.Client, method string, tgt *url.URL) *RetryableClient[req, resp] {
	return &RetryableClient[req, resp]{
		client: client,
		req:    defaultCreateRetryableRequestFunc(method, tgt, encodeHttpGenericRetryableRequest[req]),
		dec:    genericHttpResponseDecoder[resp],
	}
}

func defaultCreateRetryableRequestFunc[req any](method string, target *url.URL, enc EncodeRetryableRequestFunc[req]) CreateRetryableRequestFunc[req] {
	return func(ctx context.Context, request *req) (*retryablehttp.Request, error) {
		r, err := retryablehttp.NewRequest(method, target.String(), nil)
		if err != nil {
			return nil, err
		}

		if err = enc(ctx, r, request); err != nil {
			return nil, err
		}

		return r, nil
	}
}

// Endpoint returns a usable Go kit endpoint that calls the remote HTTP endpoint.
func (c RetryableClient[req, resp]) Endpoint() Endpoint[req, resp] {
	return func(ctx context.Context, r *req) (*resp, error) {
		ctx, cancel := context.WithCancel(ctx)

		var (
			resp *http.Response
			err  error
		)

		req, err := c.req(ctx, r)
		if err != nil {
			cancel()
			return nil, err
		}

		resp, err = c.client.Do(req.WithContext(ctx))
		if err != nil {
			cancel()
			return nil, err
		}

		defer resp.Body.Close()
		defer cancel()

		response, err := c.dec(ctx, resp)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}
