// Tideland Go Cell Network - Behaviors - Errors
//
// Copyright (C) 2010-2014 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/goas/v3/errors"
	"github.com/tideland/gocn/v2/cells"
)

//--------------------
// CONSTANTS
//--------------------

const (
	ErrIllegalRequest = iota + 1
	ErrIllegalResponse
	ErrCannotRecover
)

var errorMessages = map[int]string{
	ErrIllegalRequest:  "cell does not understand request %q",
	ErrIllegalResponse: "response has illegal type: %v",
	ErrCannotRecover:   "cannot recover cell: %v",
}

//--------------------
// ERRORS
//--------------------

// NewIllegalRequestError returns an error describing an illegal request.
func NewIllegalRequestError(request cells.Event) error {
	return errors.New(ErrIllegalRequest, errorMessages, request.Topic())
}

// IsIllegalRequestError checks if an error shows an illegal request.
func IsIllegalRequestError(err error) bool {
	return errors.IsError(err, ErrIllegalRequest)
}

// IsIllegalResponseError checks if an error shows an illegal response.
func IsIllegalResponseError(err error) bool {
	return errors.IsError(err, ErrIllegalResponse)
}

// NewCannotRecoverError returns an error showing that a cell cannot
// recover from errors.
func NewCannotRecoverError(err interface{}) error {
	return errors.New(ErrCannotRecover, errorMessages, err)
}

// IsCannotRecoverError checks if an error shows a cell that cannot
// recover.
func IsCannotRecoverError(err error) bool {
	return errors.IsError(err, ErrCannotRecover)
}

// EOF
