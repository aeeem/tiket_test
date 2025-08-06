package helper

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const (
	UniqueViolationCode = pq.ErrorCode("23505")
)

func MakeUsecaseLevelErr(statusCode int, msg string) error {
	return errors.New(strconv.Itoa(statusCode) + "|" + msg)
}

func IsErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}

func TranslateErrorToHTTPCode(err error) (code int) {
	if errors.Is(err, ErrNotFound) {
		return fiber.StatusNotFound
	}
	if errors.Is(err, Duplicate) {
		return fiber.StatusBadRequest
	}
	if errors.Is(err, BadRequest) {
		return fiber.StatusBadRequest
	}
	return fiber.StatusInternalServerError
}

func CheckIfErrFromDbToStatusCode(err error) (errs error) {
	var pqErr *pgconn.PgError
	errors.As(err, &pqErr)
	log.Print(pqErr)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if pq.ErrorCode(pqErr.Code) == UniqueViolationCode {
		return Duplicate
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return BadRequest
	}

	return InternalServerErr
}

func JsonErrorResponseValidation(c *fiber.Ctx, errs error) (err error) {
	return c.Status(400).JSON(
		ErrorResponse{
			Message:   errs.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		},
	)
}

func JsonErrorResponse(c *fiber.Ctx, errs error) (err error) {
	//if usecase level defined dont go anywhere
	isUsecase, code, msg := CheckIfErrIsUsecaseLevel(errs)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Message:   "failed translating error, internal server error!",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}

	if isUsecase {
		return c.Status(code).JSON(ErrorResponse{
			Message:   msg,
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}

	StatusCode, err := strconv.Atoi(errs.Error())
	if err != nil {
		log.Info().Err(err).Msg("Json error response error logs")
		return c.Status(500).JSON(ErrorResponse{
			Message:   "failed translating error, internal server error!",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}

	//can be using switch case but time too short for debug
	if errors.Is(errs, ErrNotFound) {
		return c.Status(StatusCode).JSON(ErrorResponse{
			Message:   "Data not found",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	} else if errors.Is(errs, Duplicate) {
		return c.Status(StatusCode).JSON(ErrorResponse{
			Message:   "Check your request Unique field can't be duplicated",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	} else if errors.Is(errs, Unauthorized) {
		return c.Status(StatusCode).JSON(ErrorResponse{
			Message:   "Unauthorized",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	} else if errors.Is(errs, Forbidden) {
		return c.Status(StatusCode).JSON(ErrorResponse{
			Message:   "Forbidden",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	} else if errors.Is(errs, BadRequest) {
		return c.Status(StatusCode).JSON(ErrorResponse{
			Message:   "Bad Request",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	} else {
		return c.Status(StatusCode).JSON(ErrorResponse{
			Message:   "Internal Server Error please contact developer",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}
