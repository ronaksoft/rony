package errors

/*
   Creation Time: 2021 - May - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Error Codes
const (
	ErrCodeInternal         = "E00" // When Error is Unknown or it is internal and should not be exposed to the client
	ErrCodeInvalid          = "E01"
	ErrCodeUnavailable      = "E02"
	ErrCodeTooMany          = "E03"
	ErrCodeTooFew           = "E04"
	ErrCodeIncomplete       = "E05"
	ErrCodeTimeout          = "E06"
	ErrCodeAccess           = "E07"
	ErrCodeAlreadyExists    = "E08"
	ErrCodeBusy             = "E09"
	ErrCodeOutOfRange       = "E10"
	ErrCodePartiallyApplied = "E11"
	ErrCodeExpired          = "E12"
	ErrCodeNotImplemented   = "E13"
)
