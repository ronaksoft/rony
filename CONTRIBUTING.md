# How to Contribute
These are mostly guidelines, not rules. Use your best judgment, and 
feel free to propose changes to this document in a pull request.

## What should I know before I get started?

There are a good few numbers of packages. All the packages which would never need to be accessed from
outside the Rony are placed in the internal package.

#### [Tools](tools)
> For developing Rony we needed many helper functions, since we think they may be handy for others too, we
placed them in this package. However, the reason of its existence was NOT to be
a library for external use.

#### [Pools](pools)
> The auto-generated codes fully attached within Rony, we made this package also available to external
worlds. This package handle pooling of frequently used structures such as Timers, WaitGroups, Buffers ...


#### [Store](store)
> Store packages contains the required functionality for accessing local storage, backed with BadgerDB. 

#### [Registry](registry)
> Registry holds names of the constructors (i.e., Aggregates, Models, Methods, Requests, Responses ...). 
In Rony messages identified by a 64bit digit instead of long named strings. This approach off-pressure
the golang runtime GC and also we wanted the lowest possible bandwidth overhead in RPC messages.

#### [Edge](edge)
> The core package of Rony. All the other components glued together in this package.
RequestCtx, DispatchCtx, EdgeServer etc defined in this package.

#### [EdgeC](edgec)
> Edge client package contains required structures and functionalities for Rony clients to interact with
Edge servers.

## Coding Convention

Start reading Rony's code you'll get the hang of it. We optimize for readability:

* We ALWAYS put spaces after list items and method parameters ([1, 2, 3], not [1,2,3]), around operators (x += 1, not x+=1), and around hash arrows.

* This is open source software. Consider the people who will read your code, and make it look nice for them. It's sort of like driving a car: Perhaps you love doing donuts when 
  you're alone, but with passengers the goal is to make the ride as smooth as possible.
  
* GolangCI linter with the provided config `.golangci.yml` file exists. Make sure you have zero warnings before your
pull request.
  
