package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	v1 "github.com/vaibhavahuja/short-video-analytics-aggregator/internal/api/v1"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/app/handlers"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/app/models"
	http2 "github.com/vaibhavahuja/short-video-analytics-aggregator/internal/endpoints/http"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue/kafka"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository/cassandra"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	engine          *gin.Engine
	repository      repository.ShortVideoRepository
	eventConsumer   queue.ConsumerInterface
	mapperProducer  queue.ProducerInterface
	reducerConsumer queue.ConsumerInterface
}

func Init() (*App, error) {
	//initialising cassandra
	var cassConfig *repository.CassandraConfig
	utils.MarshalJsonToStruct(viper.Sub("cassandra").AllSettings(), &cassConfig)
	cassandraRepo, err := cassandra.NewCassandraRepository(cassConfig)
	if err != nil {
		log.Err(err).Msg("error while initialising cassandra")
		return nil, err
	}

	apiHandler := v1.NewVideoAggregatorHandler(cassandraRepo)
	engine := http2.GetHttpServer(*apiHandler)

	var mapperProducerConfig *queue.ProducerConfig
	utils.MarshalJsonToStruct(viper.Sub("message_queue.kafka.map-producer").AllSettings(), &mapperProducerConfig)
	mapperProducer := kafka.NewProducer(mapperProducerConfig)

	//handler for event consumer does not require database access
	var eventConsumerCfg *queue.ConsumerConfig
	utils.MarshalJsonToStruct(viper.Sub("message_queue.kafka.event-consumer").AllSettings(), &eventConsumerCfg)
	eventConsumer := kafka.NewConsumer(nil, eventConsumerCfg)

	//Initialising the aggregatorMap
	models.InitAggregatorMap()

	var reducerConsumerCfg *queue.ConsumerConfig
	utils.MarshalJsonToStruct(viper.Sub("message_queue.kafka.reduce-consumer").AllSettings(), &reducerConsumerCfg)
	reducerConsumer := kafka.NewConsumer(cassandraRepo, reducerConsumerCfg)

	return &App{
		repository:      cassandraRepo,
		engine:          engine,
		eventConsumer:   eventConsumer,
		mapperProducer:  mapperProducer,
		reducerConsumer: reducerConsumer,
	}, nil
}

func (app *App) Start() error {
	ctx := context.Background()
	//starting the event and reducer handlers - both running in separate goroutines
	eventHandler := handlers.NewEventHandler(app.mapperProducer)
	app.eventConsumer.Consume(ctx, eventHandler)

	reducerHandler := handlers.NewReducerHandler(app.repository)
	//starting a processor which reads the current time and inserts into database
	go reducerHandler.ProcessCurrentTimeData(ctx)
	app.reducerConsumer.Consume(ctx, reducerHandler)

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
