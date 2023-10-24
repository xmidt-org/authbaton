// SPDX-FileCopyrightText: {{DATE}} Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

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
