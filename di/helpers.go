package di

import (
	"github.com/scylladb/gocqlx/v2"
	"gorm.io/gorm"
)

func ProvideCqlSession(s gocqlx.Session) {
	MustProvide(func() gocqlx.Session {
		return s
	})
}

func ProvideSqlSession(s *gorm.DB) {
	MustProvide(func() *gorm.DB {
		return s
	})
}
