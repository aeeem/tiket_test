package helper

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound       = errors.New("404")
	BadRequest        = errors.New("400")
	Duplicate         = errors.New("409")
	Unauthorized      = errors.New("401")
	Forbidden         = errors.New("403")
	InternalServerErr = errors.New("500")
)

type ListResponse struct {
	Status       string      `json:"success"`
	Message      string      `json:"message,omitempty" example:"success"`
	Data         interface{} `json:"results,omitempty"`
	TotalResults int         `json:"total_results,omitempty"`
	SearchID     string      `json:"search_id" example:"uuid"`
}
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type StdResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

func JsonListResponseSuccess(c *fiber.Ctx, msg string, data interface{}) (err error) {
	return c.Status(200).JSON(
		ListResponse{
			Message: msg,
			Data:    data,
			// Timestamp: time.Now().Format(time.RFC3339),
		},
	)
}

func JsonStandardResponseUpdated(c *fiber.Ctx, msg string, data interface{}) (err error) {
	return c.Status(201).JSON(
		StdResponse{
			Message:   "data has been Updated",
			Data:      data,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	)
}

func JsonStandardResponseCreated(c *fiber.Ctx, data interface{}) (err error) {
	return c.Status(201).JSON(
		StdResponse{
			Message:   "data has been created",
			Data:      data,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	)
}

func JsonStandardResponseSuccess(c *fiber.Ctx, msg string, data interface{}) (err error) {
	return c.Status(200).JSON(
		StdResponse{
			Success:   true,
			Message:   msg,
			Data:      data,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	)
}

func JsonStandardResponseDeleted(c *fiber.Ctx) (err error) {
	return c.Status(200).JSON(
		StdResponse{
			Message:   "data has been deleted",
			Data:      nil,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	)
}

func JsonErrorResponseCustomMessage(c *fiber.Ctx, errs error, msg string) (err error) {
	//checking if error is usecase level
	err = UsecaseLevelErrHTTPRespons(c, errs)
	if err == nil {
		return
	}

	return c.Status(TranslateErrorToHTTPCode(errs)).JSON(
		StdResponse{
			Message:   msg,
			Data:      nil,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	)
}
