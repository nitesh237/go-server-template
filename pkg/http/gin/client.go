package ginhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/nitesh237/go-server-template/pkg/errors"
)

// EncodeRequestFunc encodes the passed request object into the HTTP request
// object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward EncodeRequestFunc could be something that JSON
// encodes the object directly to the request body.
type EncodeRequestFunc[req any] func(context.Context, *http.Request, *req) error

// CreateRequestFunc creates an outgoing HTTP request based on the passed
// request object. It's designed to be used in HTTP clients, for client-side
// endpoints. It's a more powerful version of EncodeRequestFunc, and can be used
// if more fine-grained control of the HTTP request is required.
type CreateRequestFunc[req any] func(context.Context, *req) (*http.Request, error)

// DecodeResponseFunc extracts a user-domain response object from an HTTP
// response object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward DecodeResponseFunc could be something that
// JSON decodes from the response body to the concrete response type.
type DecodeResponseFunc[resp any] func(context.Context, *http.Response) (response *resp, err error)

// genericHttpRequestEncoder is a transport/http.EncodeRequestFunc that
// SON-encodes any request to the request body. Primarily useful in a client.
func genericHttpRequestEncoder[req any](_ context.Context, r *http.Request, request *req) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Body = io.NopCloser(bytes.NewBuffer(b))
	return nil
}

// genericHttpResponseDecoder is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//
//	decode the specific error message from the response body.
func genericHttpResponseDecoder[resp any](_ context.Context, r *http.Response) (*resp, error) {
	if r.StatusCode != http.StatusOK {
		return nil, ErrorDecoder(r)
	}

	res := new(resp)
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ErrorDecoder(r *http.Response) error {
	dc := json.NewDecoder(r.Body)
	var errResp errors.ErrorResponse
	if err := dc.Decode(&errResp); err != nil {
		return err
	}

	errResp.ErrorType = errors.GetErrorTypeFromErrorCode(errResp.Code)
	return &errResp
}

type Client[req, resp any] struct {
	client *http.Client
	req    CreateRequestFunc[req]
	dec    DecodeResponseFunc[resp]
}

// NewClient constructs a usable Client for a single remote method.
func NewClient[req, resp any](client *http.Client, method string, tgt *url.URL) *Client[req, resp] {
	return &Client[req, resp]{
		client: client,
		req:    defaultCreateRequestFunc(method, tgt, genericHttpRequestEncoder[req]),
		dec:    genericHttpResponseDecoder[resp],
	}
}

// NewClient constructs a usable Client for a single remote method.
func NewClientWithDecorator[req, resp any](client *http.Client, method string, tgt *url.URL, httpReqDecorator ...func(httpReq *http.Request, r *req) (*http.Request, error)) *Client[req, resp] {
	return &Client[req, resp]{
		client: client,
		req:    defaultCreateRequestFunc(method, tgt, genericHttpRequestEncoder[req], httpReqDecorator...),
		dec:    genericHttpResponseDecoder[resp],
	}
}

// NewClient constructs a usable Client for a single remote method.
func NewRetryableClientNative[req, resp any](client *retryablehttp.Client, method string, tgt *url.URL) *Client[req, resp] {
	return &Client[req, resp]{
		client: client.StandardClient(),
		req:    defaultCreateRequestFunc(method, tgt, genericHttpRequestEncoder[req]),
		dec:    genericHttpResponseDecoder[resp],
	}
}

func defaultCreateRequestFunc[req any](method string, target *url.URL, enc EncodeRequestFunc[req], httpReqDecorator ...func(httpReq *http.Request, r *req) (*http.Request, error)) CreateRequestFunc[req] {
	return func(ctx context.Context, request *req) (*http.Request, error) {

		req, err := http.NewRequest(method, target.String(), nil)
		if err != nil {
			return nil, err
		}

		if len(httpReqDecorator) > 0 {
			for _, httpDec := range httpReqDecorator {
				req, err = httpDec(req, request)
				if err != nil {
					return nil, errors.Wrap(err, "unable to decorate http request")
				}
			}
		}

		if err = enc(ctx, req, request); err != nil {
			return nil, err
		}

		return req, nil
	}
}

// Endpoint returns a usable Go kit endpoint that calls the remote HTTP endpoint.
func (c Client[req, resp]) Endpoint() Endpoint[req, resp] {
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
