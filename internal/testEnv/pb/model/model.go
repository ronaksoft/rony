package model

/*
   Creation Time: 2021 - Jan - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=../../../.. --go_out=. model.proto
//go:generate protoc -I=. -I=../../../.. --gorony_out=option=cql:. model.proto
func init() {}
