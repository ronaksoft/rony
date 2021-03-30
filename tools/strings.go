package tools

import (
	"encoding/binary"
	"github.com/ronaksoft/rony/pools"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

/*
   Creation Time: 2021 - Mar - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

/*
	Strings Builder helper functions
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

/*
	String Conversion helper functions
*/

func StrToInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func StrToInt32(s string) int32 {
	v, _ := strconv.ParseInt(s, 10, 32)
	return int32(v)
}

func StrToUInt64(s string) uint64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return uint64(v)
}

func StrToUInt32(s string) uint32 {
	v, _ := strconv.ParseInt(s, 10, 32)
	return uint32(v)
}

func StrToInt(s string) int {
	v, _ := strconv.ParseInt(s, 10, 32)
	return int(v)
}

func Int64ToStr(x int64) string {
	return strconv.FormatInt(x, 10)
}

func Int32ToStr(x int32) string {
	return strconv.FormatInt(int64(x), 10)
}

func UInt64ToStr(x uint64) string {
	return strconv.FormatUint(x, 10)
}

func UInt32ToStr(x uint32) string {
	return strconv.FormatUint(uint64(x), 10)
}

func IntToStr(x int) string {
	return strconv.FormatUint(uint64(x), 10)
}

// ByteToStr converts byte slice to a string without memory allocation.
// Note it may break if string and/or slice header will change
// in the future go versions.
func ByteToStr(bts []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&bts))
}

// B2S is alias for ByteToStr
func B2S(bts []byte) string {
	return ByteToStr(bts)
}

// StrToByte converts string to a byte slice without memory allocation.
// Note it may break if string and/or slice header will change
// in the future go versions.
func StrToByte(str string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

// S2B is alias for StrToByte
func S2B(str string) []byte {
	return StrToByte(str)
}
