# Tideland Go Cell Network

## Description

The *Tideland Go Cell Network* (GOCN) is a framework for event and
behavior based applications. It provides a runtime environment
for connected cells. These receive events, process them and emit
new events to their subscribers. The way how cells process the
events is defined by behaviors implementing an interface. Some
are already included.

## Installation

```
go get github.com/tideland/gocn/v3/cells
go get github.com/tideland/gocn/v3/behaviors
go get github.com/tideland/gocn/v3/testsupport
```

## Usage

### Cells

New environments are created with

```
env := cells.NewEnvironment()
```

and cells are added with

```
env.StartCell("foo", NewFooBehavior())
```

Cells then can be subscribed with

```
env.Subscribe("foo", "bar")
```

so that events emitted by the "foo" cell during the processing of
events will be received by the "bar" cell. Each cell can have
multiple cells subscibed.

Events from the outside are emitted using

```
env.Emit("foo", myEvent)
```

or

```
env.EmitNew("foo", "myTopic", cells.PayloadValues{
        "KeyA": 12345,
        "KeyB": true,
}, myScene)
```

Behaviors have to implement the `cells.Behavior` interface. Here
the `Init()` method is called with a `cells.Context`. This can be
used inside the `ProcessEvent()` method to emit events to subscribers
or directly to other cells of the environment.

Sometimes it's needed to directly communicate with a cell to retrieve
information. In this case the method

```
response, err := env.Request("foo", "myRequest?", myPayload, myScene, myTimeout)
```

is to be used. Inside the `ProcessEvent()` of the addressed cell the
event can be used to send the response with

```
switch event.Topic() {
case "myRequest?":
        event.Respond(someIncredibleData)
case ...:
        ...
}
```

Instructions without a response are simply done by emitting an event.

[![GoDoc](https://godoc.org/github.com/tideland/gocn/v3/cells?status.svg)](https://godoc.org/github.com/tideland/gocn/v3/cells)

### Behaviors

Some behaviors are already included:

- a *broadcaster* beavior that simply emits the received events to all subscribers,
- a *collector* behavior that collects events and returns them on demand,
- a *counter* behavior that counts events based on a passed function,
- a *filter* behavior that filters and emits events based on a passed function,
- a *finite state machine* bahvior that provides a simple way for state machines,
- a *logger* behavior logging the received events,
- a *mapper* behavior mapping received events to other emitted events based on
  a passed function,
- a *round-robin* behavior emitting received events round robin to its
  subscribers,
- a *router* behavior routing events to individual subscribers based on a
  passed function and
- a *ticker* behavior emitting events based on a timer.

More to come.

[![GoDoc](https://godoc.org/github.com/tideland/gocn/v3/behaviors?status.svg)](https://godoc.org/github.com/tideland/gocn/v3/behaviors)

### Testsupport

The test support package is only needed if local changes and Go tests
are planned.

[![GoDoc](https://godoc.org/github.com/tideland/gocn/v3/behaviors?testsupport.svg)](https://godoc.org/github.com/tideland/gocn/v3/testsupport)

## Authors

- Frank Mueller - <mue@tideland.biz>

## License

*Tideland Go Cell Network* is distributed under the terms of the BSD 3-Clause license.

*And now have fun.* ;)
