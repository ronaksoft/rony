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

func New(code Code, item string) *rony.Error {
	return &rony.Error{
		Code:        string(code),
		Items:       item,
		Description: "",
	}
}

func NewF(code Code, item string, format string, args ...interface{}) *rony.Error {
	return &rony.Error{
		Code:        string(code),
		Items:       item,
		Description: fmt.Sprintf(format, args...),
	}
}

func Message(reqID uint64, errCode Code, errItem string) *rony.MessageEnvelope {
	msg := &rony.MessageEnvelope{}
	ToMessage(msg, reqID, errCode, errItem)
	return msg
}

func ToMessage(out *rony.MessageEnvelope, reqID uint64, errCode Code, errItem string) {
	errMessage := rony.PoolError.Get()
	errMessage.Code = string(errCode)
	errMessage.Items = errItem
	out.Fill(reqID, rony.C_Error, errMessage)
	rony.PoolError.Put(errMessage)
}

var (
	GenAccessErr           = genWithErrorAndItem(Access)
	GenAlreadyExistsErr    = genWithErrorAndItem(AlreadyExists)
	GenBusyErr             = genWithErrorAndItem(Busy)
	GenExpiredErr          = genWithErrorAndItem(Expired)
	GenIncompleteErr       = genWithErrorAndItem(Incomplete)
	GenInternalErr         = genWithErrorAndItem(Internal)
	GenInvalidErr          = genWithErrorAndItem(Invalid)
	GenNotImplementedErr   = genWithErrorAndItem(NotImplemented)
	GenOutOfRangeErr       = genWithErrorAndItem(OutOfRange)
	GenPartiallyAppliedErr = genWithErrorAndItem(PartiallyApplied)
	GenTimeoutErr          = genWithErrorAndItem(Timeout)
	GenTooFewErr           = genWithErrorAndItem(TooFew)
	GenTooManyErr          = genWithErrorAndItem(TooMany)
	GenUnavailableErr      = genWithErrorAndItem(Unavailable)
)

func genWithErrorAndItem(code Code) func(item string, err error) *rony.Error {
	return func(item string, err error) *rony.Error {
		if err != nil {
			return &rony.Error{
				Code:        string(code),
				Items:       item,
				Description: err.Error(),
			}
		} else {
			return &rony.Error{
				Code:        string(code),
				Items:       item,
				Description: "",
			}
		}
	}
}
