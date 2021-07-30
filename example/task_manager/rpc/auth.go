package rpc

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/example/task_manager/modules/auth"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - Jul - 30
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Auth struct {
	auth.ModuleBase
}

func (s *Auth) Register(ctx *edge.RequestCtx, req *auth.RegisterRequest, res *auth.Authorization) {
	// check if username already exists
	_, err := s.Local().User.Read(req.GetUsername(), nil)
	if err == nil {
		ctx.PushError(errors.GenAlreadyExistsErr("USER", nil))
		return
	}

	err = s.Local().User.Create(&auth.User{
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
	err = s.L().Session.Create(&auth.Session{
		ID:       sessionID,
		Username: req.GetUsername(),
	})
	if err != nil {
		ctx.PushError(errors.ErrInternalServer)
		return
	}

	res.SessionID = sessionID
}

func (s *Auth) Login(ctx *edge.RequestCtx, req *auth.LoginRequest, res *auth.Authorization) {
	// check if username already exists
	_, err := s.Local().User.Read(req.GetUsername(), nil)
	if err != nil {
		ctx.PushError(errors.GenUnavailableErr("USER", nil))
		return
	}

	sessionID := tools.RandomID(32)
	err = s.L().Session.Create(&auth.Session{
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
func (s *Auth) CheckSession(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	sessionID := in.Get("SessionID", "")
	if sessionID != "" {
		s, _ := s.Local().Session.Read(sessionID, nil)
		if s != nil {
			ctx.Set("Username", s.GetUsername())
		}
	}
}

// MustAuthorized is a middleware to make sure the following handlers are only executed if
// a valid session id was provided in the request
func (s *Auth) MustAuthorized(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	username := ctx.GetString("Username", "")
	if username == "" {
		ctx.PushError(errors.GenAccessErr("SESSION", nil))
	}
}
