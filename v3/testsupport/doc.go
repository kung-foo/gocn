// Tideland Go Cell Network - Test Support
//
// Copyright (C) 2010-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The test support package provices some simple types and
// functions supporting the tests of the Go Cell Network and
// its behaviors.
package testsupport

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/goas/v1/version"
)

//--------------------
// VERSION
//--------------------

// PackageVersion returns the version of the version package.
func PackageVersion() version.Version {
	return version.New(3, 0, 0)
}

// EOF
