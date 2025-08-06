package http

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tiket_test/provider/internal/airplane"
	airplane_stream "tiket_test/provider/internal/airplane/delivery"
	repo "tiket_test/provider/internal/airplane/repository"
	usecase "tiket_test/provider/internal/airplane/usecase"

	// internalValidator "tiket_test/provider/internal/validator"
	"time"

	logs "log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	_ "tiket_test/provider/docs"
)

var validate = validator.New()

func HttpRun(port string, timeOut time.Duration) (app *fiber.App, err error) {

	app = fiber.New(
		fiber.Config{
			IdleTimeout: timeOut,
		},
	)
	newLogger := gormLogger.New(
		logs.New(os.Stdout, "\r\n", logs.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,            // Don't include params in the SQL log
			Colorful:                  true,            // Disable color
		},
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			viper.GetString("database.host"),
			viper.GetString("database.user"),
			viper.GetString("database.pass"),
			viper.GetString("database.name"),
			viper.GetString("database.port")),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	//migration
	db.AutoMigrate(
		airplane.Airplane{},
	)
	//metrics
	prometheus := fiberprometheus.New("provider")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	//docs generator
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")

	})
	app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//redis connection
	redisConn := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	red := redisConn.Ping(ctx)
	res, err := red.Result()
	if err != nil {
		panic(err)
	}
	log.Debug().Any("red", res).Msg("ping redis")

	//readfile for mocking http response
	b, err := os.ReadFile("sample.json")
	if err != nil {
		panic(err)
	}
	//pasing file into repository for cleaner code
	repository := repo.NewAirplaneRepository(redisConn, db, b, ctx)

	//usecase
	airplaneUsecase := usecase.NewAirplaneUsecase(repository)
	go func() { airplane_stream.NewAirplaneStream(redisConn, airplaneUsecase) }()

	go func() {
		if err = app.Listen(port); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	log.Print("Gracefully shutting down...")
	_ = app.Shutdown()

	log.Print("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	redisConn.Close()
	log.Print("Fiber was successful shutdown.")

	return
}
