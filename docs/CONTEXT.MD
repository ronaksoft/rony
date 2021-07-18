# Context

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

## 1. Conn
Conn is the outermost context. It is created when a connection established between client and edge
server. It is destroyed when the connection closed. Connections could be Persistent (i.e., websocket) or
Non-Persistent (i.e., http). 

## 2. DispatchCtx
DispatchCtx is created when a message (i.e., MessageEnvelop or MessageContainer) received. It will exist until the last
request executed. In other words if the received message was `MessageEnvelope` then DispatchCtx has only one `RequestCtx`
and if the received message was `MessageContainer` then DispatchCtx has as many `RequestCtx` as the number of message envelops
it contains.

## 3. RequestCtx
RequestCtx is the latest context which actually is accessed by the developer. Each RPC handler has a reference to its
`RequestCtx`. You have access to other sub-components such as Cluster, Tunnel and Store through the `RequestCtx`.

