package msg

/*
   Creation Time: 2020 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate protoc -I=../vendor -I=.  --gogofaster_out=plugins=grpc:. msg.proto
//go:generate protoc -I=../vendor -I=. --gohelpers_out=. msg.proto
var (
	ConstructorNames map[int64]string
)

func init() {
	ConstructorNames = make(map[int64]string)
}

func ErrorMessage(out *MessageEnvelope, errCode, errItem string) {
	errMessage := PoolError.Get()
	errMessage.Code = errCode
	errMessage.Items = errItem
	ResultError(out, errMessage)
	PoolError.Put(errMessage)
	return
}
