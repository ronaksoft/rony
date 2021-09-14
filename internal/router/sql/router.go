package sqlRouter

import (
	"github.com/ronaksoft/rony/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Router struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) Set(entityID string, replicaSet uint64, replace bool) error {
	e := &Entity{
		EntityID:   entityID,
		ReplicaSet: replicaSet,
		CreatedOn:  tools.TimeUnix(),
		EditedOn:   tools.TimeUnix(),
	}

	if replace {
		return r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(e).Error
	} else {
		return r.db.Create(e).Error
	}
}

func (r *Router) Get(entityID string) (replicaSet uint64, err error) {
	var entity = &Entity{}
	err = r.db.First(entity, "entity_id = ?", entityID).Error
	return entity.ReplicaSet, err
}
