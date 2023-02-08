package httpserver

import (
	"challenge/internal/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v4"
)

var validate = validator.New()

type ErrorResponse struct {
	Error string
}

type ValidationErrorResponse struct {
	Errors map[string][]string
}

func NewServer(questionRepository repository.QuestionRepository, questionOptionRepository repository.QuestionOptionRepository) *fiber.App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Mount("/questions", NewQuestionServer(questionRepository, questionOptionRepository))

	return app
}

func validateRequest(c *fiber.Ctx, req any) (res any, valid bool) {
	// Parse body
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse{Error: err.Error()}, false
	}

	// Validare
	err := validate.Struct(req)
	if err != nil {
		errs := map[string][]string{}
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructNamespace()
			errs[fieldName] = append(errs[fieldName], err.Tag())
		}
		return ValidationErrorResponse{Errors: errs}, false
	}

	return nil, true
}

func getAuthUserId(c *fiber.Ctx) (uint, error) {
	user, exists := c.Locals("user").(*jwt.Token)
	if !exists {
		return 0, nil
	}
	claims := user.Claims.(jwt.MapClaims)
	userIdU64, err := strconv.ParseUint(claims["user_id"].(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userIdU64), nil
}
