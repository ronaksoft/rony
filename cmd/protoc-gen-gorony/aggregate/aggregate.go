package aggregate

import (
	"fmt"
	"github.com/ronaksoft/rony/tools"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2021 - Jan - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Aggregate struct {
	Type      string
	Name      string
	Table     Key
	Views     []Key
	FieldsCql map[string]string
	FieldsGo  map[string]string
	HasIndex  bool
}

type Key struct {
	aggregate string
	PKs       []string
	PKGoTypes []string
	CKs       []string
	CKGoTypes []string
	Orders    []string
}

func (k *Key) Keys() []string {
	keys := make([]string, 0, len(k.PKs)+len(k.CKs))
	keys = append(keys, k.PKs...)
	keys = append(keys, k.CKs...)
	return keys
}

func (k *Key) FuncArgs(prefix string) string {
	sb := strings.Builder{}
	cnt := 0
	for idx, kk := range k.PKs {
		if cnt != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", prefix, kk)))
		sb.WriteRune(' ')
		sb.WriteString(k.PKGoTypes[idx])
		cnt++
	}
	for idx, kk := range k.CKs {
		if cnt != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", prefix, kk)))
		sb.WriteRune(' ')
		sb.WriteString(k.CKGoTypes[idx])
		cnt++
	}
	return sb.String()
}

func (k *Key) FuncArgsPKs(prefix string) string {
	sb := strings.Builder{}
	cnt := 0
	for idx, kk := range k.CKs {
		if cnt != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", prefix, kk)))
		sb.WriteRune(' ')
		sb.WriteString(k.CKGoTypes[idx])
		cnt++
	}
	return sb.String()
}

func (k *Key) FuncArgsCKs(prefix string) string {
	sb := strings.Builder{}
	cnt := 0
	for idx, kk := range k.CKs {
		if cnt != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", prefix, kk)))
		sb.WriteRune(' ')
		sb.WriteString(k.CKGoTypes[idx])
		cnt++
	}
	return sb.String()
}

func (k *Key) DBKey(prefix string) string {
	return fmt.Sprintf("'M', C_%s, %d, %s",
		k.aggregate,
		k.Checksum(),
		k.String(prefix, ",", prefix == ""),
	)
}

func (k *Key) DBKeyPrefix(prefix string) string {
	return fmt.Sprintf("'M', C_%s, %d, %s",
		k.aggregate,
		k.Checksum(),
		k.StringPKs(prefix, ",", prefix == ""),
	)
}

func (k *Key) DBIterPrefix() string {
	return fmt.Sprintf("'M', C_%s, %d",
		k.aggregate,
		k.Checksum(),
	)
}

func (k *Key) StringPKs(prefix string, sep string, lowerCamel bool) string {
	sb := strings.Builder{}
	for idx, k := range k.PKs {
		if idx != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(prefix)
		if lowerCamel {
			sb.WriteString(tools.ToLowerCamel(k))
		} else {
			sb.WriteString(k)
		}

	}
	return sb.String()
}

func (k *Key) StringCKs(prefix string, sep string, lowerCamel bool) string {
	sb := strings.Builder{}
	for idx, k := range k.CKs {
		if idx != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(prefix)
		if lowerCamel {
			sb.WriteString(tools.ToLowerCamel(k))
		} else {
			sb.WriteString(k)
		}

	}
	return sb.String()
}

func (k *Key) ChecksumPKs(keyPrefix string, sep string, lowerCamel bool) uint32 {
	return crc32.ChecksumIEEE(tools.StrToByte(k.StringPKs(keyPrefix, sep, lowerCamel)))
}

func (k *Key) String(prefix string, sep string, lowerCamel bool) string {
	sb := strings.Builder{}
	for idx, k := range k.Keys() {
		if idx != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(prefix)
		if lowerCamel {
			sb.WriteString(tools.ToLowerCamel(k))
		} else {
			sb.WriteString(k)
		}

	}
	return sb.String()
}

func (k *Key) Checksum() uint32 {
	return crc32.ChecksumIEEE(tools.StrToByte(k.String("", ",", false)))
}
