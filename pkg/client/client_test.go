package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_AccessToken = "test-access-token"

	_LoginPath         = "/nacos/v1/auth/login"
	_ConfigurationPath = "/nacos/v1/cs/configs"
)

func defaultLoginHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonResp, _ := json.Marshal(map[string]interface{}{
		"accessToken": _AccessToken,
		"tokenTtl":    18000,
	})
	_, _ = w.Write(jsonResp)
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name         string
		loginHandler http.HandlerFunc
		expectErr    error
	}{
		{
			name: "invalid response, failed to authen",
			loginHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			expectErr: fmt.Errorf("authenticate error"),
		},
		{
			name:         "success",
			loginHandler: defaultLoginHandler,
			expectErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case _LoginPath:
					assert.Equal(t, "test-user", r.FormValue("username"))
					assert.Equal(t, "test-password", r.FormValue("password"))

					tt.loginHandler(w, r)
				default:
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()

			client, err := NewClient(&Config{
				Address:     server.URL,
				Username:    "test-user",
				Password:    "test-password",
				ContextPath: "nacos",
			})
			if tt.expectErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, _AccessToken, client.accessToken)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_GetConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		getConfigHandler http.HandlerFunc
		expectErr        error
	}{
		{
			name: "request error",
			getConfigHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectErr: fmt.Errorf("internal server error"),
		},
		{
			name: "not found",
			getConfigHandler: func(w http.ResponseWriter, r *http.Request) {
				jsonResp, _ := json.Marshal(map[string]interface{}{})
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonResp)
			},
			expectErr: fmt.Errorf("not found"),
		},
		{
			name: "success",
			getConfigHandler: func(w http.ResponseWriter, r *http.Request) {
				jsonResp, _ := json.Marshal(map[string]interface{}{
					"tenant":  "namespace",
					"group":   "GROUP",
					"dataId":  "key",
					"content": "value",
				})
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonResp)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configurationId := &ConfigurationId{
				Namespace: "namespace",
				Group:     "GROUP",
				Key:       "key",
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case _LoginPath:
					defaultLoginHandler(w, r)

				case _ConfigurationPath:
					if r.Method == "GET" {
						assert.Equal(t, configurationId.Namespace, r.URL.Query().Get("tenant"))
						assert.Equal(t, configurationId.Group, r.URL.Query().Get("group"))
						assert.Equal(t, configurationId.Key, r.URL.Query().Get("dataId"))

						tt.getConfigHandler(w, r)
					}

				default:
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()

			client, err := NewClient(&Config{
				Address:     server.URL,
				ContextPath: "nacos",
			})
			assert.Nil(t, err)
			_, err = client.GetConfiguration(context.Background(), configurationId)
		})
	}
}

func TestClient_DeleteConfiguration(t *testing.T) {
	tests := []struct {
		name                string
		deleteConfigHandler http.HandlerFunc
		expectErr           error
		expectRes           bool
	}{
		{
			name: "request error",
			deleteConfigHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectErr: fmt.Errorf("internal server error"),
			expectRes: false,
		},
		{
			name: "success",
			deleteConfigHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("true"))
			},
			expectErr: nil,
			expectRes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configurationId := &ConfigurationId{
				Namespace: "namespace",
				Group:     "GROUP",
				Key:       "key",
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case _LoginPath:
					defaultLoginHandler(w, r)

				case _ConfigurationPath:
					if r.Method == "DELETE" {
						assert.Equal(t, configurationId.Namespace, r.URL.Query().Get("tenant"))
						assert.Equal(t, configurationId.Group, r.URL.Query().Get("group"))
						assert.Equal(t, configurationId.Key, r.URL.Query().Get("dataId"))

						tt.deleteConfigHandler(w, r)
					}

				default:
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()

			client, err := NewClient(&Config{
				Address:     server.URL,
				ContextPath: "nacos",
			})
			assert.Nil(t, err)
			res, err := client.DeleteConfiguration(context.Background(), configurationId)
			if tt.expectErr != nil {
				assert.NotNil(t, err)
				assert.False(t, res)
			} else {
				assert.Equal(t, tt.expectRes, res)
			}
		})
	}
}

func TestClient_PublishConfiguration(t *testing.T) {
	tests := []struct {
		name                 string
		publishConfigHandler http.HandlerFunc
		expectErr            error
	}{
		{
			name: "request error",
			publishConfigHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectErr: fmt.Errorf("internal server error"),
		},
		{
			name: "success",
			publishConfigHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("true"))
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configuration := &Configuration{
				Namespace:   "namespace",
				Group:       "GROUP",
				Key:         "key",
				Value:       "value",
				Description: "description",
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case _LoginPath:
					defaultLoginHandler(w, r)

				case _ConfigurationPath:
					if r.Method == "POST" {
						assert.Equal(t, configuration.Namespace, r.URL.Query().Get("tenant"))
						assert.Equal(t, configuration.Group, r.URL.Query().Get("group"))
						assert.Equal(t, configuration.Key, r.URL.Query().Get("dataId"))

						assert.Equal(t, configuration.Value, r.FormValue("content"))
						assert.Equal(t, configuration.Description, r.FormValue("desc"))

						tt.publishConfigHandler(w, r)
					}

				default:
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()

			client, err := NewClient(&Config{
				Address:     server.URL,
				ContextPath: "nacos",
			})
			assert.Nil(t, err)
			err = client.PublishConfiguration(context.Background(), configuration)
			if tt.expectErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
