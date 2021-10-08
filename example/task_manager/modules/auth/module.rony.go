package auth

import (
	store "github.com/ronaksoft/rony/store"
)

type ModuleBase struct {
	local LocalRepos
}

func New(store *store.Store) ModuleBase {
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
	User    *UserLocalRepo
	Session *SessionLocalRepo
}

func newLocalRepos(s *store.Store) LocalRepos {
	return LocalRepos{
		User:    NewUserLocalRepo(s),
		Session: NewSessionLocalRepo(s),
	}
}
