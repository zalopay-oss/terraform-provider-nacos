package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type requestOption struct {
	form  *url.Values
	query *url.Values
}

type requestOptionFn func(*requestOption) error

const (
	defaultPOSTContentType = "application/x-www-form-urlencoded"
	accessTokenQueryName   = "accessToken"
)

func updateValues(name string, v *url.Values, kv ...string) error {
	if len(kv)%2 == 1 {
		return fmt.Errorf("%s: odd argument count", name)
	}

	for i := 0; i < len(kv); i += 2 {
		v.Add(kv[i], kv[i+1])
	}
	return nil
}

func withQuery(kv ...string) requestOptionFn {
	return func(rOpts *requestOption) error {
		if rOpts.query == nil {
			rOpts.query = &url.Values{}
		}

		return updateValues("query string", rOpts.query, kv...)
	}
}

func withAuthentication(token *cString) requestOptionFn {
	return func(rOpts *requestOption) error {
		if rOpts.query == nil {
			rOpts.query = &url.Values{}
		}

		return updateValues("query string", rOpts.query, accessTokenQueryName, token.value())
	}
}

func withForm(kv ...string) requestOptionFn {
	return func(rOpts *requestOption) error {
		if rOpts.form == nil {
			rOpts.form = &url.Values{}
		}

		return updateValues("form data", rOpts.form, kv...)
	}
}

func newRequest(ctx context.Context, method, url string, opts ...requestOptionFn) (*http.Request, error) {
	var (
		err  error
		req  *http.Request
		body io.Reader
	)

	rOpt := &requestOption{}
	for _, opt := range opts {
		err = opt(rOpt)
		if err != nil {
			return nil, err
		}
	}

	if rOpt.form != nil {
		body = strings.NewReader(rOpt.form.Encode())
	}
	req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if rOpt.form != nil {
		req.Header.Set("Content-Type", defaultPOSTContentType)
	}

	if rOpt.query != nil {
		req.URL.RawQuery = rOpt.query.Encode()
	}

	return req, nil
}

func isTokenExpiredError(err error) bool {
	if err == nil {
		return false
	}
	re := regexp.MustCompile(`request error status_code = (\d*), body = (.*)`)
	matches := re.FindAllStringSubmatch(err.Error(), -1)
	if matches == nil || len(matches[0]) != 3 {
		return false
	}
	statusCode, body := matches[0][1], matches[0][2]
	return statusCode == "403" && strings.Contains(body, `"message":"token expired!"`)
}

func sendRequest(req *http.Request, result interface{}) error {
	var err error
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do req = %v: %v", *req, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"request error status_code = %v, body = %v",
			resp.StatusCode, string(body))
	}

	if len(body) == 0 {
		return nil
	}
	if err = json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response body = %s: %v", string(body), err)
	}
	return nil
}

// cString is a concurrent safe string
type cString struct {
	mux sync.RWMutex
	v   string
}

func (ts *cString) value() string {
	ts.mux.RLock()
	defer ts.mux.RUnlock()
	return ts.v
}

func (ts *cString) set(s string) {
	ts.mux.Lock()
	defer ts.mux.Unlock()
	ts.v = s
}
