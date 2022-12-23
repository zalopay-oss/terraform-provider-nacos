package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Address     string
	Username    string
	Password    string
	ContextPath string
}

type Client struct {
	user        *loginParams
	baseURL     string
	accessToken *cString
}

const (
	DefaultContextPath = "nacos"
	ShowAll            = "all"

	LoginPath         = "auth/login"
	ConfigurationPath = "cs/configs"
)

func NewClient(cfg *Config) (*Client, error) {
	contextPath := DefaultContextPath
	if cfg.ContextPath != "" {
		contextPath = cfg.ContextPath
	}

	client := &Client{
		user: &loginParams{
			Username: cfg.Username,
			Password: cfg.Password,
		},
		baseURL:     fmt.Sprintf("%s/%s/v1/", cfg.Address, contextPath),
		accessToken: &cString{},
	}

	if err := client.login(); err != nil {
		log.Printf("[ERROR] failed to authenticate client: %+v\n", err)
		return nil, fmt.Errorf("authenticate error: %+v\n", err)
	}

	return client, nil
}

func (c *Client) login() error {
	var resp loginResponse
	err := c.request(
		context.Background(), http.MethodPost, c.baseURL+LoginPath, &resp,
		withForm(
			"username", c.user.Username,
			"password", c.user.Password))
	if err != nil {
		return err
	}

	c.accessToken.set(resp.AccessToken)
	return nil
}

func (c *Client) request(ctx context.Context, method, url string, result interface{}, opts ...requestOptionFn) error {
	_request := func() error {
		req, err := newRequest(ctx, method, url, opts...)
		if err != nil {
			return fmt.Errorf("failed to create new request: %v", err)
		}

		err = sendRequest(req, result)
		if err != nil {
			return fmt.Errorf("failed to send request = %v: %v", *req, err)
		}

		return err
	}

	err := _request()
	if isTokenExpiredError(err) {
		if loginErr := c.login(); loginErr != nil {
			return fmt.Errorf("token expired %s, re-login attempt failed: err = %w ", err, loginErr)
		}
		err = _request()
	}

	return err
}

func (c *Client) GetConfiguration(ctx context.Context, params *ConfigurationId) (*Configuration, error) {
	var resp Configuration
	err := c.request(
		ctx, http.MethodGet, c.baseURL+ConfigurationPath, &resp,
		withAuthentication(c.accessToken),
		withQuery(
			"tenant", params.Namespace,
			"group", params.Group,
			"dataId", params.Key,
			"show", ShowAll))
	if err != nil {
		return nil, fmt.Errorf("get configuration error: %v", err)
	}
	if resp == (Configuration{}) {
		log.Printf("[WARN] not found configration=%+v\n", params)
		return nil, fmt.Errorf("not found configuration=%+v", *params)
	}

	return &resp, nil
}

func (c *Client) PublishConfiguration(ctx context.Context, params *Configuration) error {
	var resp bool
	err := c.request(
		ctx, http.MethodPost, c.baseURL+ConfigurationPath, &resp,
		withAuthentication(c.accessToken),
		withForm(
			"tenant", params.Namespace,
			"group", params.Group,
			"dataId", params.Key,
			"content", params.Value,
			"desc", params.Description))
	if err != nil {
		return fmt.Errorf("publish configuration error: %+v", err)
	}

	return nil
}

func (c *Client) DeleteConfiguration(ctx context.Context, params *ConfigurationId) (bool, error) {
	var resp bool
	err := c.request(
		ctx, http.MethodDelete, c.baseURL+ConfigurationPath, &resp,
		withAuthentication(c.accessToken),
		withQuery(
			"tenant", params.Namespace,
			"group", params.Group,
			"dataId", params.Key))
	if err != nil {
		return false, fmt.Errorf("delete configuration error: %v", err)
	}

	return true, nil
}
