package rpc

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/example/task_manager/modules/task"
	"github.com/ronaksoft/rony/store"
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

type TaskManager struct {
	task.ModuleBase
}

func (m *TaskManager) Create(ctx *edge.RequestCtx, req *task.CreateRequest, res *task.TaskView) *rony.Error {
	username := ctx.GetString("Username", "")
	newTask := &task.Task{
		ID:       tools.RandomInt64(0),
		Title:    req.GetTitle(),
		TODOs:    req.GetTODOs(),
		DueDate:  req.GetDueDate(),
		Username: username,
	}
	err := m.Local().Task.Create(newTask)
	if err != nil {
		return errors.ErrInternalServer
	}

	res.Title = newTask.GetTitle()
	res.ID = newTask.GetID()
	res.DueDate = newTask.GetDueDate()

	return nil
}

func (m *TaskManager) Get(ctx *edge.RequestCtx, req *task.GetRequest, res *task.TaskView) *rony.Error {
	t, err := m.Local().Task.Read(req.GetTaskID(), nil)
	if err != nil {
		return errors.GenUnavailableErr("TASK", err)
	}
	res.Title = t.GetTitle()
	res.ID = t.GetID()
	res.DueDate = t.GetDueDate()
	res.TODOs = t.GetTODOs()

	return nil
}

func (m *TaskManager) Delete(ctx *edge.RequestCtx, req *task.DeleteRequest, res *task.Bool) *rony.Error {
	err := m.Local().Task.Delete(req.GetTaskID())
	if err != nil {
		return errors.GenUnavailableErr("TASK", err)
	}
	res.Result = true

	return nil
}

func (m *TaskManager) List(ctx *edge.RequestCtx, req *task.ListRequest, res *task.TaskViewMany) *rony.Error {
	username := ctx.GetString("Username", "")
	tasks, err := m.Local().
		Task.
		ListByUsername(
			username,
			store.NewListOption().
				SetSkip(req.GetOffset()).
				SetLimit(req.GetLimit()),
			nil,
		)
	if err != nil {
		return errors.GenUnavailableErr("TASK", err)
	}
	for _, t := range tasks {
		res.Tasks = append(res.Tasks, &task.TaskView{
			ID:      t.GetID(),
			Title:   t.GetTitle(),
			TODOs:   t.GetTODOs(),
			DueDate: t.GetDueDate(),
		})
	}

	return nil
}
