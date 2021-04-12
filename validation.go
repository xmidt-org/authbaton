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
	"fmt"
	"net"

	"github.com/xmidt-org/themis/config"
)

// sentinel validation errors
var (
	errServerConfigMissing          = errors.New("Missing server config")
	errServerConfigUnmarshalFailure = errors.New("Failed to unmarshal server config")
	errServerAddressBadFormat       = errors.New("Server address format must be [host]:[port]")
	errServerAddressNonLoopback     = errors.New("Server address must be a loopback address")
)

type server struct {
	Address string
}

type serverValidationError struct {
	Key string
	Err error
}

func (s serverValidationError) Error() string {
	return fmt.Sprintf("%s: %v", s.Key, s.Err.Error())
}

func (s serverValidationError) Unwrap() error {
	return s.Err
}

type serverValidator struct {
	// Key is the full path to the configuration value containing
	// the server config.
	Key string
}

// Validate ensures the server's config is valid before initializing its HTTP server.
func (v serverValidator) Validate(Unmarshaller config.Unmarshaller) error {
	if !Unmarshaller.IsSet(v.Key) {
		return serverValidationError{Err: errServerConfigMissing, Key: v.Key}
	}
	var s server
	err := Unmarshaller.UnmarshalKey(v.Key, &s)
	if err != nil {
		return serverValidationError{Err: fmt.Errorf("%w: %v", errServerConfigUnmarshalFailure, err), Key: v.Key}
	}
	err = isLoopbackAddress(s.Address)
	if err != nil {
		return serverValidationError{Err: err, Key: v.Key}
	}
	return nil
}

// isLoopAddress takes an address of the form [host]:[port] and reports
// an error if the host does not refer to a loopback ip address.
// Examples of loopback addresses include 'localhost:8080' and '127.0.0.1:8080'
func isLoopbackAddress(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("%w: %v", errServerAddressBadFormat, err)
	}
	if host == "localhost" {
		return nil
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return nil
	}
	return errServerAddressNonLoopback
}
