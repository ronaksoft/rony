package rony

/*
   Creation Time: 2021 - Jul - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Tunnel provides the communication channel between edge servers. Tunnel is similar to gateway.Gateway in functionalities.
// However Tunnel is optimized for inter-communication between edge servers, and Gateway is optimized for client-server communications.
type Tunnel interface {
	Start()
	Run()
	Shutdown()
	Addr() []string
}
