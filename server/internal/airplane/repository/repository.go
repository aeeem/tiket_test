package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"tiket_test/server/internal/airplane"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type airplaneRepository struct {
	RED *redis.Client
	Ctx context.Context
}

func NewRepository(RED *redis.Client, ctx context.Context) airplane.AirplaneRepository {
	return &airplaneRepository{
		RED: RED,
		Ctx: ctx,
	}
}

func (r airplaneRepository) ACK(msgID string) (err error) {
	_, err = r.RED.XAck(r.Ctx, "search", "workers", msgID).Result()
	if err != nil {
		return
	}
	return
}
func (r airplaneRepository) GetAirplanes(search string) (result []airplane.Airplane, err error) {
	searchToINT, err := strconv.Atoi(search)
	if err != nil {
		return
	}
	searchToINT += 1
	res1, err := r.RED.XGroupCreate(r.Ctx, "search-result", search+"-mygroup", search).Result()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Print("error dari sini creategroup")
		return
	}
	log.Print(res1)
	res2, err := r.RED.XReadGroup(r.Ctx, &redis.XReadGroupArgs{
		Group:    search + "-mygroup",
		Consumer: "myconsumer",
		Block:    time.Second * 30,
		Streams:  []string{"search-result", ">"},
	}).Result()
	log.Print(res2)
	if err != nil {
		log.Print("error dari sini xreadgroup")
		log.Print(err)
		return
	}
	log.Print("sampe 3")
	for _, stream := range res2 {
		for _, msg := range stream.Messages {
			RedRes := airplane.AirplaneRedisResponse{}
			b, _ := json.Marshal(msg.Values)
			json.Unmarshal(b, &RedRes)
			temp := airplane.Airplane{
				ID:            RedRes.ID,
				Airline:       RedRes.Airline,
				FlightNumber:  RedRes.FlightNumber,
				From:          RedRes.From,
				To:            RedRes.To,
				DepartureTime: RedRes.DepartureTime,
				ArrivalTime:   RedRes.ArrivalTime,
				Currency:      RedRes.Currency,
			}
			temp.Price, err = strconv.Atoi(RedRes.Price)
			if err != nil {
				log.Err(err)
			}
			temp.Available, err = strconv.ParseBool(RedRes.Available)
			result = append(result, temp)
		}
	}

	log.Print(res2)
	return
}

func (r airplaneRepository) AddSearchQueue(AirplaneRequest airplane.AirplaneRequest) (msgID string, err error) {
	toMap := map[string]interface{}{}
	b, err := json.Marshal(AirplaneRequest)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &toMap)
	redARgs := redis.XAddArgs{

		MaxLen: 1000,
		Approx: true,
		Stream: "search-query",
		Values: toMap,
	}
	msgID, err = r.RED.XAdd(r.Ctx, &redARgs).Result()
	if err != nil {
		log.Debug().Err(err).Msg("error while adding search queue")
	}
	return
}
