package errors

import (
	"fmt"
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2021 - May - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func New(code, item string) *rony.Error {
	return &rony.Error{
		Code:        code,
		Items:       item,
		Description: "",
	}
}

func NewF(code, item string, format string, args ...interface{}) *rony.Error {
	return &rony.Error{
		Code:        code,
		Items:       item,
		Description: fmt.Sprintf(format, args...),
	}
}

func Message(reqID uint64, errCode, errItem string) *rony.MessageEnvelope {
	msg := &rony.MessageEnvelope{}
	ToMessage(msg, reqID, errCode, errItem)
	return msg
}

func ToMessage(out *rony.MessageEnvelope, reqID uint64, errCode, errItem string) {
	errMessage := rony.PoolError.Get()
	errMessage.Code = errCode
	errMessage.Items = errItem
	out.Fill(reqID, rony.C_Error, errMessage)
	rony.PoolError.Put(errMessage)
}

var (
	ErrAccess           = genWithErrorAndItem(ErrCodeAccess)
	ErrAlreadyExists    = genWithErrorAndItem(ErrCodeAlreadyExists)
	ErrBusy             = genWithErrorAndItem(ErrCodeBusy)
	ErrExpired          = genWithErrorAndItem(ErrCodeExpired)
	ErrIncomplete       = genWithErrorAndItem(ErrCodeIncomplete)
	ErrInternal         = genWithErrorAndItem(ErrCodeInternal)
	ErrInvalid          = genWithErrorAndItem(ErrCodeInvalid)
	ErrNotImplemented   = genWithErrorAndItem(ErrCodeNotImplemented)
	ErrOutOfRange       = genWithErrorAndItem(ErrCodeOutOfRange)
	ErrPartiallyApplied = genWithErrorAndItem(ErrCodePartiallyApplied)
	ErrTimeout          = genWithErrorAndItem(ErrCodeTimeout)
	ErrTooFew           = genWithErrorAndItem(ErrCodeTooFew)
	ErrTooMany          = genWithErrorAndItem(ErrCodeTooMany)
	ErrUnavailable      = genWithErrorAndItem(ErrCodeUnavailable)
)

func genWithErrorAndItem(code string) func(item string, err error) *rony.Error {
	return func(item string, err error) *rony.Error {
		if err != nil {
			return &rony.Error{
				Code:        code,
				Items:       item,
				Description: err.Error(),
			}
		} else {
			return &rony.Error{
				Code:        code,
				Items:       item,
				Description: "",
			}
		}
	}
}
