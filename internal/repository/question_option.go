package repository

import (
	"challenge/internal/entity"
	"challenge/pkg/gormprovider"
	"context"
)

type QuestionOptionRepository interface {
	gormprovider.Repository
	BulkCreateQuestionOptions(ctx context.Context, questionId uint, questionOptions []entity.QuestionOption) error
	BulkReplaceQuestionOptions(ctx context.Context, questionId uint, questionOptions []entity.QuestionOption) error
}

func NewQuestionOptionRepository(provider *gormprovider.SQLiteProvider) *questionOptionRepository {
	return &questionOptionRepository{provider.NewRepository("question_options")}
}

type questionOptionRepository struct {
	gormprovider.Repository
}

func (r *questionOptionRepository) BulkCreateQuestionOptions(ctx context.Context, questionId uint, questionOptions []entity.QuestionOption) error {
	for i := range questionOptions {
		questionOptions[i].QuestionId = questionId
	}
	return r.NewQuery(ctx).Create(&questionOptions).Error
}

func (r *questionOptionRepository) BulkReplaceQuestionOptions(ctx context.Context, questionId uint, questionOptions []entity.QuestionOption) error {
	return r.RunInTransaction(ctx, func(txCtx context.Context) error {
		err := r.NewQuery(ctx).Where("question_id", questionId).Delete(&entity.QuestionOption{}).Error
		if err != nil {
			return err
		}

		return r.BulkCreateQuestionOptions(ctx, questionId, questionOptions)
	})
}
