package gormprovider

import "gorm.io/gorm"

type Option interface {
	Apply(db *gorm.DB) *gorm.DB
}

func ApplyOptions(qry *gorm.DB, opts ...Option) *gorm.DB {
	for _, opt := range opts {
		qry = opt.Apply(qry)
	}
	return qry
}

type PreloadOption string

func (o PreloadOption) Apply(db *gorm.DB) *gorm.DB {
	return db.Preload(string(o))
}
