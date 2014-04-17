// Tideland Go Cell Network - Cells - Errors
//
// Copyright (C) 2010-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/goas/v3/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	ErrCellInit = iota
	ErrDuplicateId
	ErrInvalidId
	ErrEventRecovering
	ErrRecoveredTooOften
	ErrNoTopic
	ErrEncoding
	ErrDecoding
	ErrNoRequest
	ErrInactive
	ErrStopping
	ErrTimeout
)

var errorMessages = map[int]string{
	ErrCellInit:          "cell %q cannot initialize: %v",
	ErrDuplicateId:       "cell with id %q already exists",
	ErrInvalidId:         "cell with id %q does not exist",
	ErrEventRecovering:   "cell cannot recover after error %v: %v",
	ErrRecoveredTooOften: "cell needs too much recoverings, last error: %v",
	ErrNoTopic:           "event has no topic",
	ErrEncoding:          "error during payload encoding: %v",
	ErrDecoding:          "error during payload decoding: %v",
	ErrNoRequest:         "cannot respond, event is no request",
	ErrInactive:          "cell %q is inactive",
	ErrStopping:          "%s is stopping",
	ErrTimeout:           "operation needed too long with %v",
}

//--------------------
// ERROR CHECKING
//--------------------

// IsCellInitError checks if an error is a cell init error.
func IsCellInitError(err error) bool {
	return errors.IsError(err, ErrCellInit)
}

// IsDuplicateIdError checks if an error is a cell already exists error.
func IsDuplicateIdError(err error) bool {
	return errors.IsError(err, ErrDuplicateId)
}

// IsInvalidIdError checks if an error is a cell does not exist error.
func IsInvalidIdError(err error) bool {
	return errors.IsError(err, ErrInvalidId)
}

// IsEventRecoveringError checks if an error is an error recovering error.
func IsEventRecoveringError(err error) bool {
	return errors.IsError(err, ErrEventRecovering)
}

// IsRecoveredTooOftenError checks if an error is an illegal query error.
func IsRecoveredTooOftenError(err error) bool {
	return errors.IsError(err, ErrRecoveredTooOften)
}

// IsNoTopicError checks if an error shows that an event has no topic..
func IsNoTopicError(err error) bool {
	return errors.IsError(err, ErrNoTopic)
}

// IsEncodingError checks if an error is a payload encoding error.
func IsEncodingError(err error) bool {
	return errors.IsError(err, ErrEncoding)
}

// IsDecodingError checks if an error is a payload decoding error.
func IsDecodingError(err error) bool {
	return errors.IsError(err, ErrDecoding)
}

// IsNoRequestError checks if an error signals that an event is no request.
func IsNoRequestError(err error) bool {
	return errors.IsError(err, ErrNoRequest)
}

// IsInactiveError checks if an error is a cell inactive error.
func IsInactiveError(err error) bool {
	return errors.IsError(err, ErrInactive)
}

// IsStoppingError checks if the error shows a stopping entity.
func IsStoppingError(err error) bool {
	return errors.IsError(err, ErrStopping)
}

// IsTimeoutError checks if an error is a timeout error.
func IsTimeoutError(err error) bool {
	return errors.IsError(err, ErrTimeout)
}

//--------------------
// HELPERS
//--------------------

// newError returns the error according to code and arguments.
func newError(code int, args ...interface{}) error {
	return errors.New(code, errorMessages[code], args...)
}

// annotateError annotates the passed error with the given code and arguments.
func annotateError(err error, code int, args ...interface{}) error {
	return errors.Annotate(err, code, errorMessages[code], args...)
}

// EOF
