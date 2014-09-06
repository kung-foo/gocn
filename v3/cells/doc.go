// Tideland Go Cell Network - Cells
//
// Copyright (C) 2010-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// A framework for event and behavior based applications.
//
// Cell behaviors are defined based on an interface and can be added
// to an envrionment. Here they are running as concurrent cells that
// can be networked and communicate via events. Several useful behaviors
// are provided with the behaviors package.
package cells

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
