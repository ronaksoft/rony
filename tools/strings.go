package tools

import (
	"encoding/binary"
	"github.com/ronaksoft/rony/pools"
	"strings"
)

/*
   Creation Time: 2021 - Mar - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func AppendStrInt(sb *strings.Builder, x int) {
	b := pools.TinyBytes.GetLen(8)
	binary.BigEndian.PutUint64(b, uint64(x))
	sb.Write(b)
	pools.TinyBytes.Put(b)
	return
}

func AppendStrUInt(sb *strings.Builder, x uint) {
	b := pools.TinyBytes.GetLen(8)
	binary.BigEndian.PutUint64(b, uint64(x))
	sb.Write(b)
	pools.TinyBytes.Put(b)
	return
}

func AppendStrInt64(sb *strings.Builder, x int64) {
	b := pools.TinyBytes.GetLen(8)
	binary.BigEndian.PutUint64(b, uint64(x))
	sb.Write(b)
	pools.TinyBytes.Put(b)
	return
}

func AppendStrUInt64(sb *strings.Builder, x uint64) {
	b := pools.TinyBytes.GetLen(8)
	binary.BigEndian.PutUint64(b, x)
	sb.Write(b)
	pools.TinyBytes.Put(b)
	return
}

func AppendStrInt32(sb *strings.Builder, x int32) {
	b := pools.TinyBytes.GetLen(4)
	binary.BigEndian.PutUint32(b, uint32(x))
	sb.Write(b)
	pools.TinyBytes.Put(b)
	return
}

func AppendStrUInt32(sb *strings.Builder, x uint32) {
	b := pools.TinyBytes.GetLen(4)
	binary.BigEndian.PutUint32(b, x)
	sb.Write(b)
	pools.TinyBytes.Put(b)
	return
}
