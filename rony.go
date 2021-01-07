package rony

/*
   Creation Time: 2021 - Jan - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Conn defines the Connection interface
type Conn interface {
	ConnID() uint64
	ClientIP() string
	SendBinary(streamID int64, data []byte) error
	Persistent() bool
	Get(key string) interface{}
	Set(key string, val interface{})
}
