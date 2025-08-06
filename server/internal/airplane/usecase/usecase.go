package usecase

import (
	"tiket_test/server/internal/airplane"
)

type usecase struct {
	airplaneRepository airplane.AirplaneRepository
}

func Newusecase(repository airplane.AirplaneRepository) airplane.AirplaneUsecase {
	return &usecase{
		airplaneRepository: repository,
	}
}

func (u usecase) GetAirplanes(search string) (res []airplane.Airplane, err error) {
	res, err = u.airplaneRepository.GetAirplanes(search)
	if err != nil {
		return
	}
	err = u.airplaneRepository.ACK(search)
	if err != nil {
		return
	}
	return
}

func (u usecase) AddSearchQueue(AirplaneRequest airplane.AirplaneRequest) (msgID string, err error) {
	msgID, err = u.airplaneRepository.AddSearchQueue(AirplaneRequest)
	return
}
