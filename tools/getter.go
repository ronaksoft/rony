package tools

/*
   Creation Time: 2021 - Jul - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func GetString(v interface{}, defaultValue string) string {
	s, ok := v.(string)
	if !ok {
		return defaultValue
	}
	return s
}

func GetInt64(v interface{}, defaultValue int64) int64 {
	s, ok := v.(int64)
	if !ok {
		return defaultValue
	}
	return s
}

func GetUint64(v interface{}, defaultValue uint64) uint64 {
	s, ok := v.(uint64)
	if !ok {
		return defaultValue
	}
	return s
}

func GetInt32(v interface{}, defaultValue int32) int32 {
	s, ok := v.(int32)
	if !ok {
		return defaultValue
	}
	return s
}

func GetUint32(v interface{}, defaultValue uint32) uint32 {
	s, ok := v.(uint32)
	if !ok {
		return defaultValue
	}
	return s
}

func GetBytes(v interface{}, defaultValue []byte) []byte {
	s, ok := v.([]byte)
	if !ok {
		return defaultValue
	}
	return s
}
