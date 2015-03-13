# Tideland Go Cell Network

Status: *Work in Progress*

## Motivation

The *Tideland Go Cell Network* tries to simulate one aspect of the real world: The
recursive patterns that everything is networked in some way and reacts on events.
Nations, companies, human beings, and the cells in our brain. And non of these
cells does know everything or is responsible to handle everything. Instead they
all react based on their behavior and their individual knowledge and environment.
Sometimes they only consume events to process them, sometimes they also emit
events to those cells which they are networked with. So depending on cell behaviors,
network layout, and events very simple actions of only one cell are in the same
way possible like huge storms of working cells. Surely more likely are different
patterns of event processing inside this network.

A simple and powerful environment to create these networks is the *Tideland Go
Cell Netwrok*. Its major parts are *environments* which provide a closed network,
*cells* which provide the runtime for individual behaviors, and these *behaviors*
processing the events. Beside some supporting types `Environment` and `Cell`
are provided by the package,there are also some predefined behaviors. But most
behaviors have to be implemented by the user of the *Tideland Go Cell Netwrok*
as well as he has to define the network of the cells.

### Running a Cell Network

#### Environment

A network of cells is running inside a closed environment. The responsible type
is represented by the [`Environment`](https://godoc.org/github.com/tideland/gocn/v3/cells#Environment)
interface. A new instance is created by calling

```
env := cells.NewEnvironment()
```

This way the environment is started with default settings. Additionally one or more
options can be set with `cells.NewEnvironment(<option>, ...)`. Those options are

* `cells.Id(id string) Option` sets the environment ID, otherwise it's a UUID.
* `cells.QueueFactory(queueFactory cells.EventQueueFactory) Option` sets a factory
  creating `cells.EventQueue`s. This way own event queue implementations can be
  used for the communication between the cells. Default is a local queue.
* `cells.LocalQueueFactory(size int) Option` sets the usage of local event queues
  with an initial buffer size. The default size is 10 and it cannot be smaller
  than 5. Each cell has its own event queue.

Stopping it is later be done by calling

```
err := env.Stop()
```

It ensures the proper stopping of all cells and is also the runtime finalizer for the
environment. So it will be called automatically when the garbage collector cleans the
environment.