package helper

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func CheckIfErrIsUsecaseLevel(err error) (isUsecase bool, errcode int, msg string) {
	//split between error message and error code
	splt := strings.Split(err.Error(), "|")
	log.Info().Msg(splt[0])
	//not usecase level error
	if len(splt) != 2 {
		isUsecase = false

		return
	}

	isUsecase = true
	errcode, err = strconv.Atoi(splt[0])
	if err != nil {
		log.Info().Any("splt", splt[0]).Msg("log")
		isUsecase = false
		errcode = 500
		msg = "failed translating error at usecase level, internal tiket_test/server error!" + splt[0]
		return
	}
	msg = splt[1]
	return

}

func UsecaseLevelErrHTTPRespons(c *fiber.Ctx, errs error) (err error) {
	isUsecase, statusCode, msg := CheckIfErrIsUsecaseLevel(errs)
	if !isUsecase {
		return c.Status(statusCode).JSON(
			StdResponse{
				Message:   msg,
				Data:      nil,
				Timestamp: time.Now().Format(time.RFC3339),
			},
		)

	}
	return errs
}
