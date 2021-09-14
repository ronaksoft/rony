package sqlRouter

type Entity struct {
	EntityID   string `gorm:"primaryKey"`
	ReplicaSet uint64
	CreatedOn  int64
	EditedOn   int64
}
