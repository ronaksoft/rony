package errors

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/msg"
)

/*
   Creation Time: 2019 - Nov - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Error implements error interface and is used to convert errors to msg.MessageEnvelope
type Error struct {
	Code string
	Item string
}

// NewError instantiate a new Error and return the pointer to it
func NewError(code, item string) *Error {
	return &Error{
		Code: code,
		Item: item,
	}
}

// ToMessageEnvelope fill the input with right values
func (e *Error) ToMessageEnvelope(m *msg.MessageEnvelope) {
	msg.ErrorMessage(m, e.Code, e.Item)
}

// Error
func (e *Error) Error() string {
	return fmt.Sprintf("%s:%s", e.Code, e.Item)
}

// Error Codes
const (
	ErrCodeInternal         = "E00"
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
)

// Error Items
const (
	ErrItemApi = "API"
)
