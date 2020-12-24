package tcp

/*
   Creation Time: 2020 - Dec - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func (g *Gateway) TotalConnections() int {
	return g.totalConnections()
}
