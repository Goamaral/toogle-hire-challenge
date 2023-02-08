package httpserver

import (
	"challenge/internal/entity"
	"challenge/internal/repository"
	"challenge/pkg/env"
	"challenge/pkg/gormprovider"
	"context"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/rs/zerolog/log"
)

type QuestionServer struct {
	questionRepository       repository.QuestionRepository
	questionOptionRepository repository.QuestionOptionRepository
}

func NewQuestionServer(questionRepository repository.QuestionRepository, questionOptionRepository repository.QuestionOptionRepository) *fiber.App {
	jwtAuth := jwtware.New(jwtware.Config{SigningKey: []byte(env.GetOrDefault("JWT_SIGNING_KEY", "secret"))})
	server := &QuestionServer{questionRepository, questionOptionRepository}
	app := fiber.New()
	app.Get("/", server.ListQuestions)
	app.Post("/", jwtAuth, server.CreateQuestion)
	app.Put("/:id", jwtAuth, server.UpdateQuestion)
	app.Delete("/:id", jwtAuth, server.DeleteQuestion)

	return app
}

type ListQuestionsRequest struct {
	LastId   *uint `json:"lastId"`
	PageSize uint  `json:"pageSize" validate:"max=1000"`
	AuthorId *uint `json:"authorId"`
}

func (s QuestionServer) ListQuestions(c *fiber.Ctx) error {
	var req ListQuestionsRequest

	// Validate and parse request
	if c.Request().Header.ContentLength() > 0 {
		errRes, valid := validateRequest(c, &req)
		if !valid {
			return c.Status(fiber.StatusBadRequest).JSON(errRes)
		}
	}

	// Build question filter
	questionFilter := repository.QuestionFilter{}
	if req.AuthorId != nil {
		questionFilter.AuthorId = *req.AuthorId
	}

	// Get questions
	questions, err := s.questionRepository.ListQuestions(
		c.UserContext(),
		req.PageSize,
		req.LastId,
		gormprovider.PreloadOption("QuestionOptions"),
		questionFilter,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list questions")
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Internal error"})
	}

	return c.JSON(questions)
}

func (s QuestionServer) CreateQuestion(c *fiber.Ctx) error {
	// Get authenticated user id - author
	authorId, err := getAuthUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid JWT claims"})
	}

	// Validate and parse request
	var question entity.Question
	errRes, valid := validateRequest(c, &question)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}
	question.AuthorId = authorId

	err = s.questionRepository.RunInTransaction(c.UserContext(), func(txCtx context.Context) error {
		// Create question
		err = s.questionRepository.CreateQuestion(c.UserContext(), &question)
		if err != nil {
			return err
		}

		// Create question options
		return s.questionOptionRepository.BulkCreateQuestionOptions(c.UserContext(), question.Id, question.QuestionOptions)
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create questions")
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Internal error"})
	}

	return c.JSON(question)
}

func (s QuestionServer) UpdateQuestion(c *fiber.Ctx) error {
	// Get question id
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "id is invalid"})
	}

	// Get authenticated user id - author
	authorId, err := getAuthUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid JWT claims"})
	}

	// Validate and parse request
	var questionUpdate entity.Question
	errRes, valid := validateRequest(c, &questionUpdate)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}
	questionUpdate.Id = uint(id)

	// Check if question author is the auth user
	question, err := s.questionRepository.GetQuestion(c.UserContext(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Not found"})
	}
	if question.AuthorId != authorId {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Unauthorized"})
	}

	err = s.questionRepository.RunInTransaction(c.UserContext(), func(txCtx context.Context) error {
		// Update question
		err = s.questionRepository.UpdateQuestion(c.UserContext(), uint(id), &questionUpdate)
		if err != nil {
			return err
		}

		return s.questionOptionRepository.BulkReplaceQuestionOptions(c.UserContext(), uint(id), questionUpdate.QuestionOptions)
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update question")
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Internal error"})
	}

	return c.JSON(questionUpdate)
}

func (s QuestionServer) DeleteQuestion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "id is invalid"})
	}

	err = s.questionRepository.DeleteQuestion(c.UserContext(), uint(id))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete question")
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Internal error"})
	}

	return nil
}
