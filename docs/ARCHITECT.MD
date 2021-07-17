# Architect

Rony provides facilities to implement clustered services quickly. When you write your service
using Rony framework, you write the minimum amount of boiler-place code. Most of the code generated
by Rony.

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