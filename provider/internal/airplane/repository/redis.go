package repository

import (
	"encoding/json"
	"strings"
	"tiket_test/provider/internal/airplane"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (r *airplaneRepository) ACK(msgID []string) (err error) {
	_, err = r.redis.XAck(r.Ctx, "search-query", "workers", msgID...).Result()

	return
}

func (r *airplaneRepository) Subscribe() (search []airplane.AirplaneRequest, msgID []string, err error) {

	err = r.redis.XGroupCreateMkStream(r.Ctx, "search-query", "workers", "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Err(err)
		return
	}
	res, err := r.redis.XReadGroup(r.Ctx, &redis.XReadGroupArgs{
		Group:    "workers",
		Consumer: "worker-1",
		Streams:  []string{"search-query", ">"},
	}).Result()
	if err != nil {
		log.Print(err)

		return
	}
	for _, stream := range res {
		for _, msg := range stream.Messages {
			var data airplane.AirplaneRequest
			b, err := json.Marshal(msg.Values)
			if err != nil {
				log.Print(err)
				return []airplane.AirplaneRequest{}, []string{}, err
			}
			msgID = append(msgID, msg.ID)
			json.Unmarshal(b, &data)
			data.MessageID = msg.ID
			search = append(search, data)
		}
	}
	return
}

func (r *airplaneRepository) Publish(id string, airplane airplane.Airplane) (msgID string, err error) {
	log.Print(id)
	toMap := map[string]interface{}{}
	b, err := json.Marshal(airplane)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &toMap)

	redARgs := redis.XAddArgs{
		ID:     id,
		Stream: "search-result",
		Values: toMap,
	}
	msgID, err = r.redis.XAdd(r.Ctx, &redARgs).Result()
	if err != nil {
		log.Debug().Err(err).Msg("error while adding search queue")
	}
	log.Print("masuk redis")
	log.Print(msgID)
	return
}
