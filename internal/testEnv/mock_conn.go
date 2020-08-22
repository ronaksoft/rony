package testEnv

import "git.ronaksoft.com/ronak/rony"

/*
   Creation Time: 2020 - Apr - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type MockGatewayConn struct{}

func (m MockGatewayConn) GetAuthID() int64 {
	return 0
}

func (m MockGatewayConn) SetAuthID(int64) {
	return
}

func (m MockGatewayConn) Push(e *rony.MessageEnvelope) {

}

func (m MockGatewayConn) Pop() *rony.MessageEnvelope {
	return nil
}

func (m MockGatewayConn) GetConnID() uint64 {
	return 0
}

func (m MockGatewayConn) GetClientIP() string {
	return ""
}

func (m MockGatewayConn) SendBinary(streamID int64, data []byte) error {
	return nil
}

func (m MockGatewayConn) Flush() {
	return
}

func (m MockGatewayConn) Persistent() bool {
	return true
}
