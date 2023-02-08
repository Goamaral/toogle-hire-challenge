package entity

type Question struct {
	Id              uint             `json:"id" gorm:"primaryKey"`
	Body            string           `json:"body" validate:"required"`
	QuestionOptions []QuestionOption `json:"options" validate:"required"`
	AuthorId        uint             `json:"-"`
}

type QuestionOption struct {
	Id         uint   `json:"-" gorm:"primaryKey"`
	Body       string `json:"body" validate:"required"`
	Correct    *bool  `json:"correct" validate:"required"` // Because of validations, this bool has to be a pointer
	QuestionId uint   `json:"-"`
}
