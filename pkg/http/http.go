package http

import (
	"crypto/tls"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/nitesh237/go-server-template/pkg/cfg"
	"github.com/nitesh237/go-server-template/pkg/errors"
	"github.com/nitesh237/go-server-template/pkg/log"
)

const (
	HTTP_GET  = "GET"
	HTTP_POST = "POST"
)

// NewHttpClient creates a generic http client from the config
func NewHttpClient(httpConf *cfg.HttpClient) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   httpConf.Transport.DialContext.Timeout,
				KeepAlive: httpConf.Transport.DialContext.KeepAlive,
			}).DialContext,
			MaxIdleConns:        httpConf.Transport.MaxIdleConns,
			IdleConnTimeout:     httpConf.Transport.IdleConnTimeout,
			TLSHandshakeTimeout: httpConf.Transport.TLSHandshakeTimeout,
			MaxIdleConnsPerHost: httpConf.Transport.MaxIdleConnsPerHost,
			MaxConnsPerHost:     httpConf.Transport.MaxConnsPerHost,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: httpConf.Transport.InsecureSkipVerify,
			},
		},
		Timeout: httpConf.Timeout,
	}
}

// NewRetryableHttpClient creates a retryable http client from the config
func NewRetryableHttpClient(httpConf *cfg.HttpClient, lg log.Logger) (*retryablehttp.Client, error) {
	httpClient := NewHttpClient(httpConf)
	switch {
	case httpConf.RetryParams.ExponentialBackOff != nil:
		retryablehttp.NewClient()
		return &retryablehttp.Client{
			HTTPClient:   httpClient,
			Logger:       &retryablehttpLeveledLogger{lg: lg},
			RetryWaitMin: getTimeDurationOrDefault(httpConf.RetryParams.ExponentialBackOff.BaseInterval, time.Second),
			RetryWaitMax: getTimeDurationOrDefault(httpConf.RetryParams.ExponentialBackOff.MaxInterval, 10*time.Minute),
			RetryMax:     int(httpConf.RetryParams.ExponentialBackOff.MaxAttempts),
			CheckRetry:   retryablehttp.ErrorPropagatedRetryPolicy,
			ErrorHandler: retryablehttp.PassthroughErrorHandler,
			Backoff:      exponentialBackoffBackoff(httpConf.RetryParams.ExponentialBackOff.BaseInterval, httpConf.RetryParams.ExponentialBackOff.BackoffCoefficient),
		}, nil
	case httpConf.RetryParams.RegularInterval != nil:
		return &retryablehttp.Client{
			HTTPClient:   httpClient,
			Logger:       &retryablehttpLeveledLogger{lg: lg},
			RetryWaitMin: getTimeDurationOrDefault(httpConf.RetryParams.RegularInterval.Interval, time.Second),
			RetryWaitMax: getTimeDurationOrDefault(httpConf.RetryParams.RegularInterval.Interval, time.Minute),
			RetryMax:     int(httpConf.RetryParams.RegularInterval.MaxAttempts),
			CheckRetry:   retryablehttp.ErrorPropagatedRetryPolicy,
			ErrorHandler: retryablehttp.PassthroughErrorHandler,
			Backoff:      regularIntervalBackoff(httpConf.RetryParams.RegularInterval.Interval),
		}, nil
	default:
		return nil, errors.Wrap(errors.ErrInvalidArgument, "unknown retry policy")
	}

}

// CopyURL copies the base url and replaces the path
func CopyURL(base *url.URL, path string) *url.URL {
	n := *base
	n.Path = path
	return &n
}

// DeepCopyURL creates a deep copy of the url
func DeepCopyURL(uri *url.URL) *url.URL {
	n := *uri
	return &n
}

func regularIntervalBackoff(interval time.Duration) retryablehttp.Backoff {
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
				if s, ok := resp.Header["Retry-After"]; ok {
					if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
						return time.Second * time.Duration(sleep)
					}
				}
			}
		}

		return interval
	}
}

func exponentialBackoffBackoff(baseInterval time.Duration, exp float64) retryablehttp.Backoff {
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
				if s, ok := resp.Header["Retry-After"]; ok {
					if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
						return time.Second * time.Duration(sleep)
					}
				}
			}
		}

		return getMinDuration(time.Duration(float64(baseInterval)*math.Pow(exp, float64(attemptNum))), max)
	}
}

func getMinDuration(a, b time.Duration) time.Duration {
	time.ParseDuration("%ds")
	if a > b {
		return b
	}
	return a
}

func getTimeDurationOrDefault(d time.Duration, defaultVal time.Duration) time.Duration {
	if d == 0 {
		return defaultVal
	}

	return d
}
