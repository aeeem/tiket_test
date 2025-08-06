package repository

import (
	"context"
	"encoding/json"
	"log"
	"tiket_test/provider/internal/airplane"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type airplaneRepository struct {
	redis        *redis.Client
	Db           *gorm.DB
	Ctx          context.Context
	JSONResponse []byte
}

func NewAirplaneRepository(redis *redis.Client, Db *gorm.DB, jsonResponse []byte, Ctx context.Context) airplane.AirplaneRepository {
	seeds := []airplane.Airplane{}
	json.Unmarshal(jsonResponse, &seeds)

	for _, seed := range seeds {
		err := Db.Model(&airplane.Airplane{}).Create(&seed).Error
		if err != nil {
			log.Print(err)
		}
	}
	return &airplaneRepository{
		JSONResponse: jsonResponse,
		Ctx:          Ctx,
		Db:           Db,
		redis:        redis,
	}
}

// simulating request to db on other service using sql for easier logic
func (r *airplaneRepository) Fetch(search airplane.AirplaneRequest) (res []airplane.Airplane, err error) {
	err = r.Db.Where("airplanes.from = ?", search.From).
		Where("airplanes.to = ?", search.To).
		Where("airplanes.departure_time::date = ?", search.Date).
		Where("airplanes.total_available_passanger >= ?", search.Passengers).
		Order("airplanes.id asc").Find(&res).Error
	log.Print(err)

	if err != nil {
		log.Print(err)
		return
	}
	log.Print(res)
	return
}
