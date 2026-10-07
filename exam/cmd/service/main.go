// TODO: убрать мусор отсюда сделать нормальную инициализацию
package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	dbuslib "github.com/WantBeASleep/med_ml_lib/dbus"
	grpclib "github.com/WantBeASleep/med_ml_lib/grpc"
	observerdbuslib "github.com/WantBeASleep/med_ml_lib/observer/dbus"
	observergrpclib "github.com/WantBeASleep/med_ml_lib/observer/grpc"
	loglib "github.com/WantBeASleep/med_ml_lib/observer/log"

	"exam/internal/config"

	"github.com/ilyakaznacheev/cleanenv"

	"exam/internal/repository"

	services "exam/internal/services"

	pb "exam/internal/generated/grpc/service"

	grpchandler "exam/internal/server"

	ktprocessedsubscriber "exam/internal/dbus/consumers/ktprocessed"
	ktuploadsubscriber "exam/internal/dbus/consumers/ktupload"
	mriprocessedsubscriber "exam/internal/dbus/consumers/mriprocessed"
	mriuploadsubscriber "exam/internal/dbus/consumers/mriupload"

	dbusproducers "exam/internal/dbus/producers"
	ktprocessed "exam/internal/generated/dbus/consume/ktprocessed"
	ktupload "exam/internal/generated/dbus/consume/ktupload"
	mriprocessed "exam/internal/generated/dbus/consume/mriprocessed"
	mriupload "exam/internal/generated/dbus/consume/mriupload"

	ktpreparedpb "exam/internal/generated/dbus/produce/ktprepared"
	mricompletepb "exam/internal/generated/dbus/produce/mricomplete"
	mrisplittedpb "exam/internal/generated/dbus/produce/mrisplitted"

	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
)

const (
	successExitCode = 0
	failExitCode    = 1
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	loglib.InitLogger(loglib.WithEnv())

	cfg := config.Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("init config", "err", err)
		return failExitCode

	}

	db, err := sqlx.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		slog.Error("init db", "err", err)
		return failExitCode
	}
	defer db.Close()

	client, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Secure: false,
		Creds:  credentials.NewStaticV4(cfg.S3.Access_Token, cfg.S3.Secret_Token, ""),
	})
	if err != nil {
		slog.Error("init s3", "err", err)
		return failExitCode
	}

	if err := db.Ping(); err != nil {
		slog.Error("ping db", "err", err)
		return failExitCode
	}

	producer, err := sarama.NewSyncProducer(cfg.Broker.Addrs, nil)
	if err != nil {
		slog.Error("init sarama producer", "err", err)
		return failExitCode
	}

	producerMriSplitted := dbuslib.NewProducer[*mrisplittedpb.MriSplitted](
		producer,
		"mrisplitted",
		dbuslib.WithProducerMiddlewares[*mrisplittedpb.MriSplitted](
			observerdbuslib.CrossEventProduce,
			observerdbuslib.LogEventProduce,
		),
	)

	producerKtiPrepared := dbuslib.NewProducer[*ktpreparedpb.KtPrepared](
		producer,
		"ktprepared",
		dbuslib.WithProducerMiddlewares[*ktpreparedpb.KtPrepared](
			observerdbuslib.CrossEventProduce,
			observerdbuslib.LogEventProduce,
		),
	)

	producerMriComplete := dbuslib.NewProducer[*mricompletepb.MriComplete](
		producer,
		"mricomplete",
		dbuslib.WithProducerMiddlewares[*mricompletepb.MriComplete](
			observerdbuslib.CrossEventProduce,
			observerdbuslib.LogEventProduce,
		),
	)

	dbusAdapter := dbusproducers.New(producerMriSplitted, producerMriComplete, producerKtiPrepared)

	dao := repository.NewRepository(db, client, "mri")

	services := services.New(
		dao,
		dbusAdapter,
	)

	handler := grpchandler.New(services)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpclib.PanicRecover,
			observergrpclib.CrossServerCall,
			observergrpclib.LogServerCall,
		),
	)
	pb.RegisterExamSrvServer(server, handler)

	// dbus - операции с топиками redpanda
	mriuploadSubscriber := mriuploadsubscriber.New(services)
	mriprocessedSubscriber := mriprocessedsubscriber.New(services)
	ktuploadSubscriber := ktuploadsubscriber.New(services)
	ktprocessedSubscriber := ktprocessedsubscriber.New(services)

	mriUploadHandler := dbuslib.NewGroupSubscriber(
		"mriupload",
		cfg.Broker.Addrs,
		"mriupload",
		mriuploadSubscriber,
		dbuslib.WithSubscriberMiddlewares[*mriupload.MriUpload](
			observerdbuslib.CrossEventConsume,
			observerdbuslib.LogEventConsume,
		),
	)

	mriprocessedHandler := dbuslib.NewGroupSubscriber(
		"mriprocessed",
		cfg.Broker.Addrs,
		"mriprocessed",
		mriprocessedSubscriber,
		dbuslib.WithSubscriberMiddlewares[*mriprocessed.MriProcessed](
			observerdbuslib.CrossEventConsume,
			observerdbuslib.LogEventConsume,
		),
	)

	ktUploadHandler := dbuslib.NewGroupSubscriber(
		"ktupload",
		cfg.Broker.Addrs,
		"ktupload",
		ktuploadSubscriber,
		dbuslib.WithSubscriberMiddlewares[*ktupload.KtUpload](
			observerdbuslib.CrossEventConsume,
			observerdbuslib.LogEventConsume,
		),
	)

	ktProcessedHandler := dbuslib.NewGroupSubscriber(
		"ktprocessed",
		cfg.Broker.Addrs,
		"ktprocessed",
		ktprocessedSubscriber,
		dbuslib.WithSubscriberMiddlewares[*ktprocessed.KtProcessed](
			observerdbuslib.CrossEventConsume,
			observerdbuslib.LogEventConsume,
		),
	)

	lis, err := net.Listen("tcp", cfg.App.Url)
	if err != nil {
		slog.Error("take port", "err", err)
		return failExitCode
	}

	close := make(chan struct{})
	// ЛЮТОЕ MVP
	slog.Info("start serve", slog.String("app url", cfg.App.Url))
	go func() {
		if err := server.Serve(lis); err != nil {
			slog.Error("take port", "err", err)
			panic("serve grpc")
		}
		close <- struct{}{}
	}()
	go func() {
		// пока без DI
		if err := mriUploadHandler.Start(context.Background()); err != nil {
			slog.Error("start mriupload handler", "err", err)
		}
	}()
	go func() {
		if err := mriprocessedHandler.Start(context.Background()); err != nil {
			slog.Error("start mriprocessedHandler handler", "err", err)
		}
	}()
	go func() {
		// пока без DI
		if err := ktUploadHandler.Start(context.Background()); err != nil {
			slog.Error("start ktupload handler", "err", err)
		}
	}()
	go func() {
		if err := ktProcessedHandler.Start(context.Background()); err != nil {
			slog.Error("start ktprocessed handler", "err", err)
		}
	}()

	<-close

	return successExitCode
}
