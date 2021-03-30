package tools

import (
	"bytes"
	"reflect"
)

/*
   Creation Time: 2021 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// DeleteItemFromSlice deletes item from slice. 'slice' MUST be a slice otherwise panics.
// 'index' MUST be valid otherwise panics.
func DeleteItemFromSlice(slice interface{}, index int) {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	vLength := v.Len()
	if v.Kind() != reflect.Slice {
		panic("slice is not valid")
	}
	if index >= vLength || index < 0 {
		panic("invalid index")
	}
	switch vLength {
	case 1:
		v.SetLen(0)
	default:
		v.Index(index).Set(v.Index(v.Len() - 1))
		v.SetLen(vLength - 1)
	}
}

// SliceInt64Diff returns a - b and cb will be called on each found difference.
func SliceInt64Diff(a, b []int64, cb func(int64)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}

// SliceUint64Diff returns a - b and cb will be called on each found difference.
func SliceUint64Diff(a, b []uint64, cb func(uint64)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}

// SliceInt32Diff returns a - b and cb will be called on each found difference.
func SliceInt32Diff(a, b []int32, cb func(int32)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}

// SliceUint32Diff returns a - b and cb will be called on each found difference.
func SliceUint32Diff(a, b []uint32, cb func(uint32)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}

// SliceIntDiff returns a - b and cb will be called on each found difference.
func SliceIntDiff(a, b []int, cb func(int)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}

// SliceUintDiff returns a - b and cb will be called on each found difference.
func SliceUintDiff(a, b []uint, cb func(uint)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}

// SliceInt64Diff returns a - b and cb will be called on each found difference.
func SliceStringDiff(a, b []string, cb func(string)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}

// SliceBytesDiff returns a - b and cb will be called on each found difference.
func SliceBytesDiff(a, b [][]byte, cb func([]byte)) {
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if bytes.Equal(a[i], b[j]) {
				found = true
				break
			}
		}
		if !found && cb != nil {
			cb(a[i])
		}
	}
}
