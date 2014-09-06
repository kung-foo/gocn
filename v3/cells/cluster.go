// Tideland Go Cell Network - Cells - Cluster
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
	"fmt"
	"strings"
	"sync"

	"github.com/tideland/goas/v3/errors"
)

//--------------------
// CELL CLUSTER
//--------------------

// cluster is a map from id to cells for subscriptions
// and subscribers. It is also responsible for creating and
// stopping cells or to emit events to them.
type cluster struct {
	mux   sync.RWMutex
	cells map[string]*cell
}

// newCluster creates a new cell cluster.
func newCluster() *cluster {
	return &cluster{
		cells: make(map[string]*cell),
	}
}

// startCell starts a new cell in the cluster. It will
// only be called by the environment.
func (c *cluster) startCell(env *environment, id string, behavior Behavior) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	// Check if the id already exists.
	if _, ok := c.cells[id]; ok {
		return errors.New(ErrDuplicateId, errorMessages, id)
	}
	// Create cell.
	cell, err := newCell(env, id, behavior)
	if err != nil {
		return err
	}
	c.cells[id] = cell
	return nil
}

// stopCell stops a cell in the cluster. It will only be
// called by the environment.
func (c *cluster) stopCell(id string) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if cell, ok := c.cells[id]; ok {
		delete(c.cells, id)
		return cell.stop()
	}
	return errors.New(ErrInvalidID, errorMessages, id)
}

// cell returns the cell with the given id.
func (c *cluster) cell(id string) (*cell, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if cell, ok := c.cells[id]; ok {
		return cell, nil
	}
	return nil, errors.New(ErrInvalidID, errorMessages, id)
}

// subscribe adds cells to the cluster.
func (c *cluster) subscribe(scm *cluster) {
	c.mux.Lock()
	defer c.mux.Unlock()
	scm.mux.RLock()
	defer scm.mux.RUnlock()
	for id, cell := range scm.cells {
		c.cells[id] = cell
	}
}

// unsubscribe removes cells from the cluster.
func (c *cluster) unsubscribe(ccm *cluster) {
	c.mux.Lock()
	defer c.mux.Unlock()
	ccm.mux.RLock()
	defer ccm.mux.RUnlock()
	for id := range ccm.cells {
		delete(c.cells, id)
	}
}

// ids returns the cell ids inside the cluster.
func (c *cluster) ids() []string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	cmids := []string{}
	for id := range c.cells {
		cmids = append(cmids, id)
	}
	return cmids
}

// subscribers returns the subscriber ids of a cell.
func (c *cluster) subscribers(id string) ([]string, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if cell, ok := c.cells[id]; ok {
		return cell.subscribers.ids(), nil
	}
	return nil, errors.New(ErrInvalidID, errorMessages, id)
}

// subset returns a cell cluster containing the cells with the given ids.
func (c *cluster) subset(ids ...string) (*cluster, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	sc := newCluster()
	for _, id := range ids {
		cell, ok := c.cells[id]
		if !ok {
			return nil, errors.New(ErrInvalidID, errorMessages, id)
		}
		sc.cells[id] = cell
	}
	return sc, nil
}

// emitDirect emits an event to one cell of the cluster.
func (c *cluster) emitDirect(id string, event Event) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	cell, ok := c.cells[id]
	if !ok {
		return errors.New(ErrInvalidID, errorMessages, id)
	}
	return cell.processEvent(event)
}

// emit emits an event to all cells of the cluster.
func (c *cluster) emit(event Event) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	for _, cell := range c.cells {
		if err := cell.processEvent(event); err != nil {
			return err
		}
	}
	return nil
}

// stop stops all cell of the cluster.
func (c *cluster) stop() {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, cell := range c.cells {
		cell.stop()
	}
}

// String returns a readable representation of the cell cluster.
func (c *cluster) String() string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return fmt.Sprintf("<cell cluster %s>", strings.Join(c.ids(), "; "))
}

// EOF
