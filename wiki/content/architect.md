---
title: "Architect"
date: 2021-07-26T14:33:20+04:30
draft: false
weight: 50
menu: "main"
tags: ["wiki", "install", "guide"]
categories: ["wiki"]

weight: 30
contentCopyright: MIT
mathjax: true
autoCollapseToc: false

---

# Big Picture
Rony provides facilities to implement clustered services quickly. When you write your service
using Rony framework, you write the minimum amount of boiler-place code. Most of the code generated
by Rony.

{{< image src="/images/rony_overview.png" alt="Architect Diagram" position="center" style="border-radius: 8px;" >}}

## Components
### Edge
The outermost layer of Rony is `edge.Server`. This is the only object which you need to
access all the capabilities of Rony. Edge servers are identifies by their unique id (i.e., replica set)
in a cluster. When you write your code you can do different actions depending on the replica-set of the
current node, find other edge servers by their replica-set or their unique id.

Edge servers are build of other sub-components. Gateway, Cluster, Tunnel and Store. In following we
explain what is their responsibility of each component.

#### 1. Gateway
Gateway is the only component of the edge which MUST be setup. You cannot have an edge server without
any gateway defined. Gateway manages the connection between `edge.Server` and external clients. Clients
could talk to the gateway by `http`, `websocket` or both.

#### 2. Cluster
Cluster is an OPTIONAL component of the edge server. Cluster manages the cluster as its name spoiled.
You set up cluster when initializing your edge server. Cluster communicates with each other by using gossip protocol.
You can get information about the current node or any other node in the cluster by accessing `Cluster` components.

#### 3. Tunnel
Tunnel is an OPTIONAL component of the edge server. Tunnel provides end-to-end communication between edge server.
Rony recommends to use `Tunnel` with the information collected from `Cluster` component.

#### 4. Store
Store is a REQUIRED component of the edge server. Store provides persistent key-value store. It is useful and could be
used to by you to extend your ability to save structured data locally. It supports TTL on keys.

---

## Context

Rony has three type contexts. They provide three layer of separation. `rony.Conn` `edge.DispatchCtx` and
`edge.RequestCtx` are the layered contexts of Rony.
```
|--- Conn ----------------------------------------|
|   |--- DispatchCtx --------------------------|  |
|   |   |--- RequestCtx -------------------|   |  |
|   |   |--- RequestCtx -------------------|   |  |
|   |   |   .                              |   |  |
|   |   |   .                              |   |  |
|   |------------------------------------------|  |
|                                                 |
|   |--- DispatchCtx --------------------------|  |
|   |   |--- RequestCtx -------------------|   |  |
|   |   |--- RequestCtx -------------------|   |  |
|   |------------------------------------------|  |
|                                                 |
|                      .                          |
|                      .                          |
|-------------------------------------------------|
```

### 1. Conn
Conn is the outermost context. It is created when a connection established between client and edge
server. It is destroyed when the connection closed. Connections could be Persistent (i.e., websocket) or
Non-Persistent (i.e., http).

### 2. DispatchCtx
DispatchCtx is created when a message (i.e., MessageEnvelop or MessageContainer) received. It will exist until the last
request executed. In other words if the received message was `MessageEnvelope` then DispatchCtx has only one `RequestCtx`
and if the received message was `MessageContainer` then DispatchCtx has as many `RequestCtx` as the number of message envelops
it contains.

### 3. RequestCtx
RequestCtx is the latest context which actually is accessed by the developer. Each RPC handler has a reference to its
`RequestCtx`. You have access to other sub-components such as Cluster, Tunnel and Store through the `RequestCtx`.

