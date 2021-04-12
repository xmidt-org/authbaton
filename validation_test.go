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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLoopbackAddress(t *testing.T) {
	tcs := []struct {
		Description string
		Input       string
		ExpectedErr error
	}{
		{
			Description: "Bad format",
			Input:       "127.0.0.1",
			ExpectedErr: errServerAddressBadFormat,
		},
		{
			Description: "Non loopback address",
			Input:       "remote-host.example.net:8090",
			ExpectedErr: errServerAddressNonLoopback,
		},
		{
			Description: "Localhost",
			Input:       "localhost:80",
		},
		{
			Description: "Loopback ip",
			Input:       "127.0.0.1:8080",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			assert := assert.New(t)
			err := isLoopbackAddress(tc.Input)
			assert.True(errors.Is(err, tc.ExpectedErr))
		})
	}
}
