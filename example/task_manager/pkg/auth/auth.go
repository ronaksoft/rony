package auth

/*
   Creation Time: 2021 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=$GOPATH/src -I=../../vendor --go_out=paths=source_relative:. auth.proto
func init() {}