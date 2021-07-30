package task

import (
	rony "github.com/ronaksoft/rony"
)

type ModuleBase struct {
	local LocalRepos
}

func New(store rony.Store) ModuleBase {
	m := ModuleBase{
		local: newLocalRepos(store),
	}
	return m
}

func (m ModuleBase) Local() LocalRepos {
	return m.local
}

func (m ModuleBase) L() LocalRepos {
	return m.local
}

type LocalRepos struct {
	Task *TaskLocalRepo
}

func newLocalRepos(s rony.Store) LocalRepos {
	return LocalRepos{
		Task: NewTaskLocalRepo(s),
	}
}
