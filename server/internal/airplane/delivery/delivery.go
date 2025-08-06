package delivery

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tiket_test/server/internal/airplane"
	"tiket_test/server/internal/helper"
	"tiket_test/server/internal/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type SearchRequest struct {
	airplane.AirplaneRequest
}

type SearchResponse struct {
	SearchID string `json:"search_id"`
	Status   string `json:"status"`
}
type DeliveryHandler struct {
	Validator *validator.XValidator
	usecase   airplane.AirplaneUsecase
}

func NewDelivery(app *fiber.App, validator *validator.XValidator, usecase airplane.AirplaneUsecase) {
	handler := DeliveryHandler{
		usecase:   usecase,
		Validator: validator,
	}
	app.Post("/api/flight/search", handler.PostAirplane)
	app.Get("/api/flight/search/:search_id/stream", handler.GetAirPlane)
}

// @Router			/api/flight/search	[post]
// @Summary			Adding search queue
// @Tags			Server
// @Accept			json
// @Param			id	body	delivery.SearchRequest	true "register search"
// @Success			200	{object}	helper.StdResponse{data=delivery.SearchResponse}
func (h *DeliveryHandler) PostAirplane(c *fiber.Ctx) (err error) {
	request := SearchRequest{}
	err = json.Unmarshal(c.Body(), &request)
	if err != nil {
		return helper.JsonErrorResponseValidation(c, err)
	}
	if errs := h.Validator.Validate(request); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		err = errors.New(strings.Join(errMsgs, "/n"))
		return helper.JsonErrorResponseValidation(c, err)
	}

	res, err := h.usecase.AddSearchQueue(request.AirplaneRequest)
	if err != nil {
		return helper.JsonErrorResponse(c, err)
	}
	return helper.JsonStandardResponseSuccess(c, "Search request submitted", SearchResponse{SearchID: res, Status: "pending"})
}

var sseMessageQueue []string

func (h *DeliveryHandler) GetAirPlane(c *fiber.Ctx) (err error) {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	searchID := c.Params("search_id")
	searchIDs := strings.Split(searchID, "-")
	log.Debug().Any("searchIDs", searchIDs).Msg("searchIDs")
	QuerySearch, err := strconv.Atoi(searchIDs[0])
	if err != nil {
		return helper.JsonErrorResponseValidation(c, errors.New("invalid search id"))
	}
	QuerySearch = QuerySearch - 1

	if err != nil {
		return helper.JsonErrorResponse(c, err)
	}
	airplaneResult := []airplane.Airplane{}
	result := helper.ListResponse{}
	result.SearchID = searchID
	result.Status = "pending"
	ResultList := []airplane.Airplane{}
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

		fmt.Println("WRITER")
		for {
			var msg string
			res, err := h.usecase.GetAirplanes(strconv.Itoa(QuerySearch))
			if err != nil {
				log.Debug().Err(err)
			}
			if len(res) == 0 && result.Status == "completed" {
				log.Print("Darisini")
				result.Data = nil
				result.TotalResults = len(ResultList)
				result.Status = "completed"
				log.Debug().Any("result", result).Msg("result")
				b, err := json.Marshal(result)
				if err != nil {
					log.Debug().Err(err).Msg("error when marshaling json")
				}
				msg = fmt.Sprintf("%s", string(b))
				fmt.Fprintf(w, "data:%s\n\n", msg)
				fmt.Println(msg)
				err = w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				break //closing after last message
			}

			airplaneResult = append(airplaneResult, res...)
			// if there are messages that have been sent to the `/publish` endpoint
			// then use these first, otherwise just send the current time
			if len(airplaneResult) > 0 {
				ResultList = append(ResultList, airplaneResult[0])
				result.Status = "completed"
				result.Data = ResultList
				log.Debug().Any("result", result).Msg("result")
				b, _ := json.Marshal(result)
				msg = fmt.Sprintf("%s", string(b))
				airplaneResult = airplaneResult[1:]
			} else {
				result.Data = nil
				result.TotalResults = len(ResultList)
				result.Status = "completed"
				log.Debug().Any("result", result).Msg("result")
				b, err := json.Marshal(result)
				if err != nil {
					log.Debug().Err(err).Msg("error when marshaling json")
				}
				msg = fmt.Sprintf("%s", string(b))
				fmt.Fprintf(w, "data:%s\n\n", msg)
				fmt.Println(msg)
				break //exit after last message

			}

			fmt.Fprintf(w, "data:%s\n\n", msg)
			fmt.Println(msg)

			err = w.Flush()
			if err != nil {
				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

				break
			}
			time.Sleep(2 * time.Second)
		}
	}))

	return nil

}
