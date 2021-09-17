/**
 * Copyright 2021 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
			f, err := onErrorHTTPResponse(onErrorHTTPResponseConfig{AuthType: tc.AuthType})
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
