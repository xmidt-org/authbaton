package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xmidt-org/bascule/basculehttp"
)

func TestOnErrorHTTPResponse(t *testing.T) {
	tcs := []struct {
		Description  string
		AuthType     string
		ErrReason    basculehttp.ErrorResponseReason
		ShouldFail   bool
		ExpectedCode int
	}{
		{
			Description: "Unsupported Type",
			AuthType:    "Digest",
			ShouldFail:  true,
		},
		{
			Description:  "Checks not found",
			AuthType:     "Basic",
			ErrReason:    basculehttp.ChecksNotFound,
			ExpectedCode: http.StatusForbidden,
		},
		{
			Description:  "Checks failed",
			AuthType:     "Basic",
			ErrReason:    basculehttp.ChecksFailed,
			ExpectedCode: http.StatusForbidden,
		},
		{
			Description:  "Parse failed",
			AuthType:     "Bearer",
			ErrReason:    basculehttp.ParseFailed,
			ExpectedCode: http.StatusForbidden,
		},
		{
			Description:  "No Authorization Header",
			AuthType:     "Bearer",
			ErrReason:    basculehttp.MissingHeader,
			ExpectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			assert := assert.New(t)
			f, err := onErrorHTTPResponse(OnErrorHTTPResponseOption{AuthType: tc.AuthType})
			if tc.ShouldFail {
				assert.NotNil(err)
			} else {
				recorder := httptest.NewRecorder()
				f(recorder, tc.ErrReason)
				assert.Equal(tc.ExpectedCode, recorder.Code)
				if tc.ExpectedCode == http.StatusUnauthorized {
					assert.Equal(tc.AuthType, recorder.Header().Get(basculehttp.AuthTypeHeaderKey))
				} else {
					assert.Empty(recorder.Header().Get(basculehttp.AuthTypeHeaderKey))
				}
			}
		})
	}
}
