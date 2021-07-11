package rpc

/*
   Creation Time: 2021 - Jul - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// remove the extra / to enable go:generate commands
///go:generate protoc -I=. -I=$GOPATH/src -I=../vendor --go_out=paths=source_relative:. sample.proto
///go:generate protoc -I=. -I=$GOPATH/src -I=../vendor --gorony_out=paths=source_relative:. sample.proto

func init() {}
