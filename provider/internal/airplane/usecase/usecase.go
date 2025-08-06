package usecase

import (
	"fmt"
	"log"
	"strings"
	"tiket_test/provider/internal/airplane"
	"time"
)

type airplaneUsecase struct {
	AirplaneRepository airplane.AirplaneRepository
}

func NewAirplaneUsecase(repository airplane.AirplaneRepository) airplane.AirplaneUsecase {
	return &airplaneUsecase{
		AirplaneRepository: repository,
	}
}

func (u airplaneUsecase) Publish(search airplane.AirplaneRequest) (err error) {
	res, err := u.AirplaneRepository.Fetch(search)
	if err != nil {
		return
	}
	log.Print(res)
	for _, v := range res {
		messageID := strings.Split(search.MessageID, "-")
		msgID := fmt.Sprintf("%s-*", messageID[0])
		_, err = u.AirplaneRepository.Publish(msgID, v)
		if err != nil {
			log.Print(err)
			return
		}
		time.Sleep(20 * time.Second)
	}
	return
}
func (u airplaneUsecase) Subscribe() (search []airplane.AirplaneRequest, err error) {
	search, msgID, err := u.AirplaneRepository.Subscribe()
	if err != nil {
		log.Print(err)

		return
	}
	err = u.AirplaneRepository.ACK(msgID)

	return
}
