package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	http2 "github.com/vaibhavahuja/short-video-analytics-aggregator/internal/endpoints/http"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue/kafka"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository/cassandra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	engine *gin.Engine

	eventConsumer   queue.ConsumerInterface
	mapperProducer  queue.ProducerInterface
	reducerConsumer queue.ConsumerInterface
}

func Init() (*App, error) {
	engine := http2.GetHttpServer()

	cassandraRepo := cassandra.NewCassandraRepository()
	//todo add configs in producer and consumer to connect to the correct instance
	mapperProducer := kafka.NewProducer()
	//handler for event consumer does not require database access
	eventConsumer := kafka.NewConsumer(nil)
	reducerConsumer := kafka.NewConsumer(cassandraRepo)

	return &App{
		engine:          engine,
		eventConsumer:   eventConsumer,
		mapperProducer:  mapperProducer,
		reducerConsumer: reducerConsumer,
	}, nil
}

func (app *App) Start() error {
	startServer(app.engine)
	return nil
}

// start http server here
func startServer(engine *gin.Engine) {
	port := viper.GetInt("server.port")
	if port == 0 {
		log.Fatal().Err(errors.New("http port not configured")).Send()
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: engine,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	//waits to quit until it receives the SIGTERM
	quitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(quitSignalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignalChannel
}
