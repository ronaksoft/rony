package errors

/*
   Creation Time: 2021 - May - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Code string

// Error Codes
const (
	Internal         Code = "E00" // When Error is Unknown or it is internal and should not be exposed to the client
	Invalid          Code = "E01"
	Unavailable      Code = "E02"
	TooMany          Code = "E03"
	TooFew           Code = "E04"
	Incomplete       Code = "E05"
	Timeout          Code = "E06"
	Access           Code = "E07"
	AlreadyExists    Code = "E08"
	Busy             Code = "E09"
	OutOfRange       Code = "E10"
	PartiallyApplied Code = "E11"
	Expired          Code = "E12"
	NotImplemented   Code = "E13"
)

func (c Code) Name() string {
	return _codeName[c]
}

var _codeName = map[Code]string{
	Internal:         "Internal",
	Invalid:          "Invalid",
	Unavailable:      "Unavailable",
	TooMany:          "TooMany",
	TooFew:           "TooFew",
	Incomplete:       "Incomplete",
	Timeout:          "Timeout",
	Access:           "Access",
	AlreadyExists:    "AlreadyExists",
	Busy:             "Busy",
	OutOfRange:       "OutOfRange",
	PartiallyApplied: "PartiallyApplied",
	Expired:          "Expired",
	NotImplemented:   "NotImplemented",
}
