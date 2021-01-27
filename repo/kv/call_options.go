package kv

/*
   Creation Time: 2021 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type ListOption struct {
	skip    int32
	limit   int32
	reverse bool
}

func NewListOption() *ListOption {
	return &ListOption{
		skip:    0,
		limit:   100,
		reverse: false,
	}
}

func (o *ListOption) SetSkip(skip int32) *ListOption {
	o.skip = skip
	return o
}

func (o *ListOption) Skip() int32 {
	return o.skip
}

func (o *ListOption) SetLimit(limit int32) *ListOption {
	if limit > 0 && limit < 1000 {
		o.limit = limit
	}
	return o
}

func (o *ListOption) Limit() int32 {
	return o.limit
}

func (o *ListOption) SetForward() *ListOption {
	o.reverse = false
	return o
}

func (o *ListOption) SetBackward() *ListOption {
	o.reverse = true
	return o
}

func (o *ListOption) Backward() bool {
	return o.reverse
}
