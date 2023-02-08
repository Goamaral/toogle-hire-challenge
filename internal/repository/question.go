package repository

import (
	"challenge/internal/entity"
	"challenge/pkg/gormprovider"
	"context"

	"gorm.io/gorm/clause"
)

type QuestionRepository interface {
	gormprovider.Repository
	ListQuestions(ctx context.Context, pageSize uint, lastId *uint, opts ...gormprovider.Option) ([]entity.Question, error)
	CreateQuestion(ctx context.Context, question *entity.Question) error
	GetQuestion(ctx context.Context, id uint) (entity.Question, error)
	UpdateQuestion(ctx context.Context, id uint, question *entity.Question) error
	DeleteQuestion(ctx context.Context, id uint) error
}

func NewQuestionRepository(provider *gormprovider.SQLiteProvider) *questionRepository {
	return &questionRepository{provider.NewRepository("questions")}
}

type questionRepository struct {
	gormprovider.Repository
}

func (r *questionRepository) ListQuestions(ctx context.Context, pageSize uint, lastId *uint, opts ...gormprovider.Option) ([]entity.Question, error) {
	if pageSize == 0 {
		pageSize = 10
	}

	qry := gormprovider.ApplyOptions(r.NewQuery(ctx), opts...).Limit(int(pageSize))
	if lastId != nil {
		qry = qry.Where("id > ?", *lastId)
	}

	var questions []entity.Question
	err := qry.Find(&questions).Error

	return questions, err
}

func (r *questionRepository) GetQuestion(ctx context.Context, id uint) (entity.Question, error) {
	var question entity.Question
	err := r.NewQuery(ctx).Where("id", id).First(&question).Error
	return question, err
}

func (r *questionRepository) CreateQuestion(ctx context.Context, question *entity.Question) error {
	return r.NewQuery(ctx).Omit(clause.Associations).Create(question).Error
}

func (r *questionRepository) UpdateQuestion(ctx context.Context, id uint, question *entity.Question) error {
	return r.NewQuery(ctx).Omit(clause.Associations).Where("id", id).Updates(&question).Error
}

func (r *questionRepository) DeleteQuestion(ctx context.Context, id uint) error {
	return r.NewQuery(ctx).Delete(&entity.Question{Id: id}).Error
}
