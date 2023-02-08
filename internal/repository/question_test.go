package repository_test

import (
	"challenge/internal/entity"
	"challenge/internal/repository"
	"challenge/pkg/gormprovider"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addQuestion(t *testing.T, sqlProvider *gormprovider.SQLiteProvider, question *entity.Question) *entity.Question {
	err := sqlProvider.DB.Table("questions").Create(question).Error
	if err != nil {
		t.Fatal(err)
	}
	return question
}

func TestQuestionRepository_ListQuestions(t *testing.T) {
	var lastId uint = 1
	var authorId uint = 1

	type Test struct {
		TestName    string
		PageSize    uint
		LastId      *uint
		Opts        []gormprovider.Option
		ExpectedIds []uint
	}
	tests := []Test{
		{
			TestName:    "Default arguments",
			ExpectedIds: []uint{1, 2},
		},
		{
			TestName:    "Set PageSize",
			PageSize:    1,
			ExpectedIds: []uint{1},
		},
		{
			TestName:    "Set LastId",
			LastId:      &lastId,
			ExpectedIds: []uint{2},
		},
		{
			TestName:    "Use QuestionFilter",
			Opts:        []gormprovider.Option{repository.QuestionFilter{AuthorId: authorId}},
			ExpectedIds: []uint{2},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			sqlProvider := gormprovider.NewTestSQLiteProvider(t, getDatabaseInitSql(t))
			addQuestion(t, sqlProvider, &entity.Question{Id: 1})
			addQuestion(t, sqlProvider, &entity.Question{Id: 2, AuthorId: authorId})

			repo := repository.NewQuestionRepository(sqlProvider)
			questions, err := repo.ListQuestions(context.Background(), test.PageSize, test.LastId, test.Opts...)
			require.NoError(t, err)
			questionIds := make([]uint, len(questions))
			for i, q := range questions {
				questionIds[i] = q.Id
			}
			assert.ElementsMatch(t, test.ExpectedIds, questionIds)
		})
	}
}

func TestQuestionRepository_CreateQuestion(t *testing.T) {
	// TODO
}

func TestQuestionRepository_GetQuestion(t *testing.T) {
	// TODO
}

func TestQuestionRepository_UpdateQuestion(t *testing.T) {
	// TODO
}

func TestQuestionRepository_DeleteQuestion(t *testing.T) {
	// TODO
}
