package main

/*
   Creation Time: 2020 - Jul - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Generator int

const (
	_ Generator = iota
	Pools
	RPC
	ModelScylla
)
