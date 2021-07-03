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
	Type       string
	Name       string
	Table      Key
	Views      []Key
	ViewParams []string
	FieldNames []string
	FieldsCql  map[string]string
	FieldsGo   map[string]string
	HasIndex   bool
}

func (m *Aggregate) FuncArgs(keyPrefix string, key Key) string {
	sb := strings.Builder{}
	for idx, k := range key.Keys() {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", keyPrefix, k)))
		sb.WriteRune(' ')
		sb.WriteString(m.FieldsGo[k])
	}
	return sb.String()
}

func (m *Aggregate) FuncArgsPKs(keyPrefix string, key Key) string {
	sb := strings.Builder{}
	for idx, k := range key.PKs {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", keyPrefix, k)))
		sb.WriteRune(' ')
		sb.WriteString(m.FieldsGo[k])
	}
	return sb.String()
}

func (m *Aggregate) FuncArgsCKs(keyPrefix string, key Key) string {
	sb := strings.Builder{}
	for idx, k := range key.CKs {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", keyPrefix, k)))
		sb.WriteRune(' ')
		sb.WriteString(m.FieldsGo[k])
	}
	return sb.String()
}

type Key struct {
	aggregate string
	PKs       []string
	CKs       []string
	Orders    []string
}

func (k *Key) Keys() []string {
	keys := make([]string, 0, len(k.PKs)+len(k.CKs))
	keys = append(keys, k.PKs...)
	keys = append(keys, k.CKs...)
	return keys
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
