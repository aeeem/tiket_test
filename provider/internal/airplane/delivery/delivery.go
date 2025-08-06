package airplane_stream

import (
	"log"
	usecase "tiket_test/provider/internal/airplane"
	"time"

	"github.com/redis/go-redis/v9"
)

type AirplaneStream struct {
	redis   *redis.Client
	usecase usecase.AirplaneUsecase
}

func NewAirplaneStream(redis *redis.Client, usecase usecase.AirplaneUsecase) {
	stream := AirplaneStream{
		redis:   redis,
		usecase: usecase,
	}

	err := stream.PubSub()
	if err != nil {
		return
	}

}

func (h AirplaneStream) PubSub() (err error) {
	for true {
		searchRequest, err := h.usecase.Subscribe()
		if err != nil {
			log.Print(err)
			return err
		}

		for _, v := range searchRequest {
			err = h.usecase.Publish(v)
			log.Print("masuk")
			if err != nil {
				log.Print(err)
				return err
			}
			time.Sleep(1 * time.Second) //simulating delay
		}
	}
	return
}
