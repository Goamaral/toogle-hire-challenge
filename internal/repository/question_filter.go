package repository

import "gorm.io/gorm"

type QuestionFilter struct {
	AuthorId uint
}

func (f QuestionFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.AuthorId != 0 {
		db.Where("questions.author_id", f.AuthorId)
	}

	return db
}
