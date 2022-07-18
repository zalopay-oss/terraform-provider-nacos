package client

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	testcases := []struct {
		name  string
		query []string
		form  []string
		err   error
	}{
		{
			name:  "invalid query",
			query: []string{"a", "b", "c"},
			err:   fmt.Errorf("odd argument count"),
		},
		{
			name: "invalid form",
			form: []string{"a", "b", "c"},
			err:  fmt.Errorf("odd argument count"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := request(
				context.Background(), "GET", "/test", nil,
				withForm(tc.form...), withQuery(tc.query...))
			assert.True(t, strings.Contains(err.Error(), tc.err.Error()))
		})
	}
}
