package http

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	airplaneHandler "tiket_test/server/internal/airplane/delivery"
	airplaneRepository "tiket_test/server/internal/airplane/repository"
	airplaneUsecase "tiket_test/server/internal/airplane/usecase"

	internalValidator "tiket_test/server/internal/validator"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/spf13/viper"

	_ "tiket_test/server/docs"
)

var validate = validator.New()

func HttpRun(port string) {

	myValidator := &internalValidator.XValidator{
		Validator: validate,
	}
	log.Print(myValidator)

	app := fiber.New()

	prometheus := fiberprometheus.New("provider")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")

	})
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	//redis connection
	redisConn := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx, cancel := context.WithCancel(context.TODO())
	red := redisConn.Ping(ctx)
	res, err := red.Result()
	if err != nil {
		panic(err)
	}
	log.Debug().Any("red", res).Msg("ping redis")
	//repository
	airplaneRepository := airplaneRepository.NewRepository(redisConn, ctx)
	airplaneUsecase := airplaneUsecase.Newusecase(airplaneRepository)
	airplaneHandler.NewDelivery(app, myValidator, airplaneUsecase)

	//gracefully shutdown
	go func() {

		if err = app.Listen(port); err != nil {
			panic(err)
		}

	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	_ = <-c
	log.Print("Gracefully shutting down...")
	err = app.Shutdown()
	if err != nil {
		log.Debug().Err(err).Msg("error while shutting down")
	}

	log.Print("Running cleanup tasks...")

	// close redis conn
	redisConn.Close()
	log.Print("Fiber was successful shutdown.")

}
