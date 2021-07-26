---
title: "Introduction"
date: 2021-07-26T14:33:20+04:30
draft: false
menu: "main"
weight: 50
tags: ["wiki"]
categories: ["wiki"]

weight: 10
contentCopyright: MIT
mathjax: true
autoCollapseToc: false

---

### What is Rony ?
Rony is a scalable and cluster-aware RPC framework. It also supports REST apis using generated codes.

### Why choose Rony?
The first question comes to our mind that why not use other web-server frameworks. They are
quite good at what they are intended to do. They all scale well, but they are more or less well suited
in the three tier architectures. We use Rony when our api-servers (i.e., we call it Edge servers) need
to communicate with each other directly, and they need to somehow sync with each other. Or pass their
clients messages to other clients connected to other Edge servers. This architecture shines when the 
number of apis are limited per each service. That reminds me microservice architectures. Yes in
micro-service architecture one challenge is how to scale up the microservice itself. Sync Rony provide
us cluster and end-to-end communication in the cluster out of the box, it is the best choice for such scenarios.
Rony builds its cluster using a gossip protocol hence we don't need to add another 3rd party service to our 
infrastructure.

Code generation is the other trained feature of Rony. We all hate boilerplate codes, and we have our own unique
technique individually or forced by our company regulates. Rony tries to understand the most from our proto files.
Rony defines a good few sets of options and generates code for us.



### Performance vs Conventions
Rony do a lot of tasks, but it should not add too much overhead computationally and memory wise. Rony architected to be
independent as much as possible to GC. Since GC is the overhead for Rony and Rony is overhead for your main business logic code.
The rule of thumb is In Rony do not use functions' arguments out of the context of that function. If you need to then try to clone
it. Clone function provided almost everywhere in Rony. It is encouraged to use sync.WaitGroup in the handlers and wait for all the 
sub-tasks before leaving rpc handler context. Rony takes care of your concurrency under the hood.




