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
	"git.tideland.biz/goas/errors"
	"git.tideland.biz/gocn/cells"
)

//--------------------
// CONSTANTS
//--------------------

const (
	ecNoError = iota
	ecIllegalRequest
	ecIllegalResponse
	ecCannotRecover

	msgIllegalRequest  = "cell does not understand request %q"
	msgIllegalResponse = "response has illegal type: %v"
	msgCannotRecover   = "cannot recover cell: %v"
)

//--------------------
// ERRORS
//--------------------

// NewIllegalRequestError returns an error describing an illegal request.
func NewIllegalRequestError(request cells.Event) error {
	return errors.New(ecIllegalRequest, msgIllegalRequest, request.Topic())
}

// IsIllegalRequestError checks if an error shows an illegal request.
func IsIllegalRequestError(err error) bool {
	return errors.IsError(err, ecIllegalRequest)
}

// IsIllegalResponseError checks if an error shows an illegal response.
func IsIllegalResponseError(err error) bool {
	return errors.IsError(err, ecIllegalResponse)
}

// NewCannotRecoverError returns an error showing that a cell cannot
// recover from errors.
func NewCannotRecoverError(err interface{}) error {
	return errors.New(ecCannotRecover, msgCannotRecover, err)
}

// IsCannotRecoverError checks if an error shows a cell that cannot
// recover.
func IsCannotRecoverError(err error) bool {
	return errors.IsError(err, ecCannotRecover)
}

// EOF
