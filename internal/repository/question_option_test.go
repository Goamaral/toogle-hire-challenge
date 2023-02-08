package repository_test

import (
	"challenge/internal/entity"
	"challenge/pkg/gormprovider"
	"testing"
)

func addQuestionOption(t *testing.T, sqlProvider *gormprovider.SQLiteProvider, questionOption *entity.QuestionOption) *entity.QuestionOption {
	err := sqlProvider.DB.Table("question_options").Create(questionOption).Error
	if err != nil {
		t.Fatal(err)
	}
	return questionOption
}

func TestQuestionOptionRepository_BulkCreateQuestionOptions(t *testing.T) {
	// TODO
}

func TestQuestionOptionRepository_BulkReplaceQuestionOptions(t *testing.T) {
	// TODO
}
