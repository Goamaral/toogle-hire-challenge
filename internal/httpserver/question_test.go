package httpserver_test

import (
	"bytes"
	"challenge/internal/entity"
	"challenge/internal/httpserver"
	"challenge/internal/repository"
	"challenge/mocks"
	"challenge/pkg/gormprovider"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListQuestions(t *testing.T) {
	var lastId uint = 1
	var authorId uint = 1
	var questionId uint = 1
	correct := true
	incorrect := false
	question := entity.Question{
		Id:   questionId,
		Body: "question",
		QuestionOptions: []entity.QuestionOption{
			{
				Body:       "question option correct",
				Correct:    &correct,
				QuestionId: questionId,
			},
			{
				Body:       "question option incorrect",
				Correct:    &incorrect,
				QuestionId: questionId,
			},
		},
		AuthorId: authorId,
	}
	expectedResponseBytes, err := json.Marshal([]entity.Question{question})
	require.NoError(t, err)

	type Test struct {
		TestName               string
		Req                    httpserver.ListQuestionsRequest
		ExpectedQuestionFilter repository.QuestionFilter
	}
	tests := []Test{
		{
			TestName: "EmptyRequest",
		},
		{
			TestName: "WithLastId",
			Req:      httpserver.ListQuestionsRequest{LastId: &lastId},
		},
		{
			TestName: "WithPageSize",
			Req:      httpserver.ListQuestionsRequest{PageSize: 20},
		},
		{
			TestName:               "WithAutorId",
			Req:                    httpserver.ListQuestionsRequest{AuthorId: &authorId},
			ExpectedQuestionFilter: repository.QuestionFilter{AuthorId: authorId},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			questionRepository := mocks.NewQuestionRepository(t)
			questionRepository.
				On(
					"ListQuestions",
					mock.Anything,
					test.Req.PageSize,
					test.Req.LastId,
					gormprovider.PreloadOption("QuestionOptions"),
					test.ExpectedQuestionFilter,
				).
				Return([]entity.Question{question}, nil)

			reqBodyBytes, err := json.Marshal(test.Req)
			require.NoError(t, err)

			server := httpserver.NewServer(questionRepository, nil)
			req := httptest.NewRequest(http.MethodGet, "/questions", bytes.NewReader(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			res, err := server.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusOK, res.StatusCode)

			resBodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, expectedResponseBytes, resBodyBytes)
		})
	}
}

func TestCreateQuestion(t *testing.T) {
	var questionId uint = 1
	var authorId uint = 1
	correct := true
	incorrect := false
	question := entity.Question{
		Id:   questionId,
		Body: "question",
		QuestionOptions: []entity.QuestionOption{
			{
				Body:    "question option correct",
				Correct: &correct,
			},
			{
				Body:    "question option incorrect",
				Correct: &incorrect,
			},
		},
	}
	successExpectedResponseBytes, err := json.Marshal(question)
	require.NoError(t, err)

	validJwt, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"user_id": strconv.FormatUint(uint64(authorId), 10)},
	).SignedString([]byte("secret"))
	require.NoError(t, err)
	validAuthHeader := fmt.Sprintf("Bearer %s", validJwt)

	invalidJwt, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"user_id": strconv.FormatUint(uint64(authorId), 10)},
	).SignedString([]byte("invalid_secret"))
	require.NoError(t, err)
	invalidAuthHeader := fmt.Sprintf("Bearer %s", invalidJwt)

	type Test struct {
		TestName                  string
		Req                       map[string]any
		ReqAuthHeader             string
		ExpectedHttpStatusCode    int
		ExpectedResponseBodyBytes []byte
	}
	tests := []Test{
		{
			TestName: "Success",
			Req: map[string]any{
				"body": question.Body,
				"options": []map[string]any{
					{"body": question.QuestionOptions[0].Body, "correct": question.QuestionOptions[0].Correct},
					{"body": question.QuestionOptions[1].Body, "correct": question.QuestionOptions[1].Correct},
				},
			},
			ReqAuthHeader:             validAuthHeader,
			ExpectedHttpStatusCode:    http.StatusOK,
			ExpectedResponseBodyBytes: successExpectedResponseBytes,
		},
		{
			TestName:                  "Unauthorized",
			ReqAuthHeader:             invalidAuthHeader,
			ExpectedHttpStatusCode:    http.StatusUnauthorized,
			ExpectedResponseBodyBytes: []byte("Invalid or expired JWT"),
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			questionRepository := mocks.NewQuestionRepository(t)
			if test.ExpectedHttpStatusCode == http.StatusOK {
				questionRepository.On("RunInTransaction", mock.Anything, mock.Anything).
					Return(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				questionRepository.On("CreateQuestion", mock.Anything, mock.Anything).
					Return(func(_ context.Context, q *entity.Question) error {
						assert.Equal(t, authorId, q.AuthorId)
						*q = question
						q.Id = questionId
						return nil
					})
			}

			questionOptionRepository := mocks.NewQuestionOptionRepository(t)
			if test.ExpectedHttpStatusCode == http.StatusOK {
				questionOptionRepository.On("BulkCreateQuestionOptions", mock.Anything, questionId, question.QuestionOptions).
					Return(func(_ context.Context, _ uint, questionOptions []entity.QuestionOption) error {
						for i := range questionOptions {
							questionOptions[i].Id = uint(i)
							questionOptions[i].QuestionId = questionId
						}
						return nil
					})
			}

			reqBodyBytes, err := json.Marshal(test.Req)
			require.NoError(t, err)

			server := httpserver.NewServer(questionRepository, questionOptionRepository)
			req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", test.ReqAuthHeader)
			res, err := server.Test(req)
			require.NoError(t, err)
			require.Equal(t, test.ExpectedHttpStatusCode, res.StatusCode)

			resBodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.ExpectedResponseBodyBytes, resBodyBytes)
		})
	}
}

func TestUpdateQuestion(t *testing.T) {
	// TODO
}
