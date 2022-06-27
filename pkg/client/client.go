package client

import (
	"context"
	"fmt"
	"log"
	"net/http"

	builder "github.com/phuc1998/http-builder"
)

type Config struct {
	Address     string
	Username    string
	Password    string
	ContextPath string
}

type Client struct {
	user        *loginParams
	httpClient  *builder.APIClient
	accessToken string
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
	httpBuilderCfg := builder.NewConfiguration().
		AddBasePath(fmt.Sprintf("%s/%s/v1/", cfg.Address, contextPath)).
		AddHTTPClient(&http.Client{})
	httpBuilderClient := builder.NewAPIClient(httpBuilderCfg)
	loginParams := &loginParams{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	client := &Client{
		user:       loginParams,
		httpClient: httpBuilderClient,
	}

	if err := client.login(); err != nil {
		log.Printf("[ERROR] failed to authenticate client: %+v\n", err)
		return nil, fmt.Errorf("authenticate error: %+v\n", err)
	}

	return client, nil
}

func (c *Client) login() error {
	var resp loginResponse

	_, err := c.httpClient.Builder(LoginPath).
		Post().
		BuildRequest(c.user).
		UseMultipartFormData().
		Call(context.Background(), &resp)
	if err != nil {
		return err
	}

	c.accessToken = resp.AccessToken
	return nil
}

func (c *Client) getAuthParams() *authParams {
	return &authParams{
		AccessToken: c.accessToken,
	}
}

func (c *Client) GetConfiguration(ctx context.Context, params *ConfigurationId) (*Configuration, error) {
	resp := &Configuration{}
	_, err := c.httpClient.Builder(ConfigurationPath).
		Get().
		BuildQuery(c.getAuthParams()).
		BuildQuery(optionalParams{Show: ShowAll}).
		BuildRequest(params).
		Call(ctx, resp)

	if err != nil {
		return nil, fmt.Errorf("get configuration error: %v", err)
	}
	if *resp == (Configuration{}) {
		log.Printf("[WARN] not found configration=%+v", params)
		return nil, fmt.Errorf("not found configuration=%+v", *params)
	}

	return resp, nil
}

func (c *Client) PublishConfiguration(ctx context.Context, params *Configuration) error {
	var resp bool
	_, err := c.httpClient.Builder(ConfigurationPath).
		Post().
		BuildQuery(c.getAuthParams()).
		BuildRequest(params).
		Call(ctx, &resp)

	if err != nil {
		return fmt.Errorf("publish configuration error: %+v", err)
	}

	return nil
}

func (c *Client) DeleteConfiguration(ctx context.Context, params *ConfigurationId) (bool, error) {
	var resp bool
	_, err := c.httpClient.Builder(ConfigurationPath).
		Delete().
		BuildQuery(c.getAuthParams()).
		BuildRequest(params).
		Call(ctx, resp)

	if err != nil {
		return false, fmt.Errorf("delete configuration error: %v", err)
	}

	return true, nil
}
