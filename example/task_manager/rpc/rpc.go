package rpc

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/example/task_manager/pkg/auth"
	"github.com/ronaksoft/rony/example/task_manager/pkg/task"
	"github.com/ronaksoft/rony/tools"
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
	userRepo    *auth.UserLocalRepo
	sessionRepo *auth.SessionLocalRepo
}

func NewAuth(s rony.Store) *Auth {
	return &Auth{
		userRepo:    auth.NewUserLocalRepo(s),
		sessionRepo: auth.NewSessionLocalRepo(s),
	}
}

func (a *Auth) Register(ctx *edge.RequestCtx, req *RegisterRequest, res *Authorization) {
	// check if username already exists
	_, err := a.userRepo.Read(req.GetUsername(), nil)
	if err == nil {
		ctx.PushError(errors.GenAlreadyExistsErr("USER", nil))
		return
	}

	err = a.userRepo.Create(&auth.User{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
	})
	if err != nil {
		ctx.PushError(errors.ErrInternalServer)
		return
	}

	sessionID := tools.RandomID(32)
	err = a.sessionRepo.Create(&auth.Session{
		ID:       sessionID,
		Username: req.GetUsername(),
	})
	if err != nil {
		ctx.PushError(errors.ErrInternalServer)
		return
	}

	res.SessionID = sessionID
}

func (a *Auth) Login(ctx *edge.RequestCtx, req *LoginRequest, res *Authorization) {
	// check if username already exists
	_, err := a.userRepo.Read(req.GetUsername(), nil)
	if err != nil {
		ctx.PushError(errors.GenUnavailableErr("USER", nil))
		return
	}

	sessionID := tools.RandomID(32)
	err = a.sessionRepo.Create(&auth.Session{
		ID:       sessionID,
		Username: req.GetUsername(),
	})
	if err != nil {
		ctx.PushError(errors.ErrInternalServer)
		return
	}

	res.SessionID = sessionID
}

// CheckSession is a middleware try to load the user info from the session id.
func (a *Auth) CheckSession(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	sessionID := in.Get("SessionID", "")
	if sessionID != "" {
		s, _ := a.sessionRepo.Read(sessionID, nil)
		if s != nil {
			ctx.Set("Username", s.GetUsername())
		}
	}
}

// MustAuthorized is a middleware to make sure the following handlers are only executed if
// a valid session id was provided in the request
func (a *Auth) MustAuthorized(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	username := ctx.GetString("Username", "")
	if username == "" {
		ctx.PushError(errors.GenAccessErr("SESSION", nil))
	}
}

type TaskManager struct{
	taskRepo *task.TaskLocalRepo
}

func NewTaskManager(s rony.Store) *TaskManager {
	return &TaskManager{
		taskRepo: task.NewTaskLocalRepo(s),
	}
}

func (tm *TaskManager) Create(ctx *edge.RequestCtx, req *CreateRequest, res *TaskView) {
	username := ctx.GetString("Username", "")
	newTask := &task.Task{
		ID:       tools.RandomInt64(0),
		Title:    req.GetTitle(),
		TODOs:    req.GetTODOs(),
		DueDate:  req.GetDueDate(),
		Username: username,
	}
	err := tm.taskRepo.Create(newTask)
	if err != nil {
		ctx.PushError(errors.ErrInternalServer)
		return
	}

	res.Title = newTask.GetTitle()
	res.ID = newTask.GetID()
	res.DueDate = newTask.GetDueDate()
}

func (tm *TaskManager) Get(ctx *edge.RequestCtx, req *GetRequest, res *TaskView) {
	t, err := tm.taskRepo.Read(req.GetTaskID(), nil)
	if err != nil {
		ctx.PushError(errors.GenUnavailableErr("TASK", err))
		return
	}
	res.Title = t.GetTitle()
	res.ID = t.GetID()
	res.DueDate = t.GetDueDate()
	res.TODOs = t.GetTODOs()
}

func (tm *TaskManager) Delete(ctx *edge.RequestCtx, req *DeleteRequest, res *Bool) {
	panic("implement me")
}

func (tm *TaskManager) List(ctx *edge.RequestCtx, req *ListRequest, res *TaskViewMany) {
	panic("implement me")
}
