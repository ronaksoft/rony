package rpc

import (
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

func (m *TaskManager) Create(ctx *edge.RequestCtx, req *task.CreateRequest, res *task.TaskView) {
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
		ctx.PushError(errors.ErrInternalServer)
		return
	}

	res.Title = newTask.GetTitle()
	res.ID = newTask.GetID()
	res.DueDate = newTask.GetDueDate()
}

func (m *TaskManager) Get(ctx *edge.RequestCtx, req *task.GetRequest, res *task.TaskView) {
	t, err := m.Local().Task.Read(req.GetTaskID(), nil)
	if err != nil {
		ctx.PushError(errors.GenUnavailableErr("TASK", err))
		return
	}
	res.Title = t.GetTitle()
	res.ID = t.GetID()
	res.DueDate = t.GetDueDate()
	res.TODOs = t.GetTODOs()
}

func (m *TaskManager) Delete(ctx *edge.RequestCtx, req *task.DeleteRequest, res *task.Bool) {
	err := m.Local().Task.Delete(req.GetTaskID())
	if err != nil {
		ctx.PushError(errors.GenUnavailableErr("TASK", err))
		return
	}
	res.Result = true

}

func (m *TaskManager) List(ctx *edge.RequestCtx, req *task.ListRequest, res *task.TaskViewMany) {
	username := ctx.GetString("Username", "")
	tasks, err := m.Local().Task.ListByUsername(username, store.NewListOption().SetSkip(req.GetOffset()).SetLimit(req.GetLimit()), nil)
	if err != nil {
		ctx.PushError(errors.GenUnavailableErr("TASK", err))
		return
	}
	for _, t := range tasks {
		res.Tasks = append(res.Tasks, &task.TaskView{
			ID:      t.GetID(),
			Title:   t.GetTitle(),
			TODOs:   t.GetTODOs(),
			DueDate: t.GetDueDate(),
		})
	}
}
