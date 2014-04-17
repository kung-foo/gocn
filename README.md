Tideland Go Cell Network
========================

Description
-----------

The *Tideland Go Cell Network* (GOCN) is a framework for event and
behavior based applications. It provides a runtime environment
for connected cells. These receive events, process them and emit
new events to their subscribers. The way how cells process the
events is defined by behaviors implementing an interface. Some
are already included.

Installation
------------

    go get github.com/tideland/gocn/v2/cells
    go get github.com/tideland/gocn/v2/behaviors
    go get github.com/tideland/gocn/v2/testsupport

Usage
-----

**Cells**

A runtime environment for cells is created with

    env := cells.NewEnvironment("myEnvironment")

Here cells can be started and stopped with

    err := env.StartCell("foo", NewFooBehavior())
    err = env.StartCell("bar", NewBarBehavior())
    err = env.StartCell("yadda", NewYaddaBehavior())

    err = env.StopCell("xyz")

The same behavior or behavior instance can be started multiple times
with different ids. Cells can be subscribed to each other with

    err = env.Subscribe("foo", "bar", "yadda")

In this case events emitted by *foo* will be receivd by *bar* and
*yadda* to be processed there. Initial events can be created and 
emitted with

    err = env.Raise("foo", "myTopic", somePayload)

A behavior has to implement methods for initialization, termination,
error recovering and naturally for the processing of events.

    func (m *MyBehavior) ProcessEvent(event Event) error {
            ...
    }

During initialization a cell context is passed to the cell. That
can be used to access the environment or emit one or more events 
to the own subscribers with

    err := ctx.RaiseAll("anotherTopc", anotherPayload)

**Behaviors**

Some behaviors are already included:

- a broadcaster beavior that simply emits the received events to all subscribers,
- a collector behavior that collects events and returns them on demand,
- a counter behavior that counts events based on a passed function,
- a filter behavior that filters and emits events based on a passed function,
- a finite state machine bahvior that provides a simple way for state machines,
- a logger behavior logging the received events,
- a mapper behavior mapping received events to other emitted events based on
  a passed function,
- a round-robin behavior emitting received events round robin to its
  subscribers,
- a router behavior routing events to individual subscribers based on a
  passed function and
- a ticker behavior emitting events based on a timer.

More to come.

**Test Support**

The test support package is only needed if local changes and Go tests
are planned.

And now have fun. ;)

Documentation
-------------

- http://godoc.org/github.com/tideland/gocn/v2/cells
- http://godoc.org/github.com/tideland/gocn/v2/behaviors
- http://godoc.org/github.com/tideland/gocn/v2/testsupport

Authors
-------

- Frank Mueller - <mue@tideland.biz>

License
-------

*Tideland Go Cell Network* is distributed under the terms of the BSD 3-Clause license.
