package bridge

/*
   Creation Time: 2019 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Bridge interface {
	GetID() string
	SendNotifier(bridgeID string, connID uint64)
	SendMessage(bridgeID string, authID int64, data []byte) bool
}
