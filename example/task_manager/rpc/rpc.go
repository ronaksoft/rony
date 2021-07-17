package rpc

import (
	"github.com/ronaksoft/rony/edge"
)

/*
   Creation Time: 2021 - Jul - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=../../.. --go_out=paths=source_relative:. auth.proto task.proto
//go:generate protoc -I=. -I=../../.. --gorony_out=paths=source_relative:. auth.proto task.proto
func init() {}

type Auth struct {

}

func (a *Auth) Register(ctx *edge.RequestCtx, req *RegisterRequest, res *Authorization) {
	panic("implement me")
}

func (a *Auth) Login(ctx *edge.RequestCtx, req *LoginRequest, res *Authorization) {
	panic("implement me")
}


type TaskManager struct {}

func (t *TaskManager) Create(ctx *edge.RequestCtx, req *CreateRequest, res *TaskView) {
	panic("implement me")
}

func (t *TaskManager) Get(ctx *edge.RequestCtx, req *GetRequest, res *TaskView) {
	panic("implement me")
}

func (t *TaskManager) Delete(ctx *edge.RequestCtx, req *DeleteRequest, res *Bool) {
	panic("implement me")
}

func (t *TaskManager) List(ctx *edge.RequestCtx, req *ListRequest, res *TaskViewMany) {
	panic("implement me")
}
