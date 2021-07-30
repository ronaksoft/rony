package task

/*
   Creation Time: 2021 - Jul - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=../../../.. --go_out=paths=source_relative:. model.proto rpc.proto
//go:generate protoc -I=. -I=../../../.. --gorony_out=paths=source_relative,rony_opt=module:. model.proto rpc.proto
func init() {}
