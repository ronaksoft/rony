package pb

/*
   Creation Time: 2020 - Mar - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate protoc -I=.  --gogofaster_out=. service.proto
//go:generate protoc -I=. --gorony_out=. service.proto
func init() {

}

var (
	ConstructorNames = map[int64]string{}
)
