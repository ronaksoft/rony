package rony

import (
	"git.ronaksoft.com/ronak/rony/tools"
	"hash/crc64"
	"reflect"
)

/*
   Creation Time: 2020 - Apr - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	crc64Table = crc64.MakeTable(crc64.ISO)
)

func ConstructorOf(x interface{}) int64 {
	return int64(crc64.Checksum(tools.StrToByte(reflect.ValueOf(x).Type().Name()), crc64Table))
}
