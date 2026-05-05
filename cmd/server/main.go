package main

import (
	g "Price/gen/bot"
	"Price/internal/config"
	"Price/internal/logger"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	h "Price/internal/delivery/grpc"
	health "Price/internal/delivery/health_check"
	repo "Price/internal/repository"
	sP "Price/internal/service/product"
	sS "Price/internal/service/subscription"
	sU "Price/internal/service/user"
	w "Price/internal/worker"

	par "Price/gen/parser"

	inf "Price/internal/infrastructure"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {

	wg := &sync.WaitGroup{}

	cfg := config.MustLoad()
	l := logger.NewLogger()

	l.Info("сервер запускается", "host", cfg.DB.Host, "port", cfg.DB.Port)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Printf("Error opening database: %v", err)
		panic(err)
	}
	defer db.Close()

	cb := inf.NewCircuitBreaker(5, 1*time.Minute)

	// infrastructure redis/kafka

	rDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",
		DB:       0})

	defer rDB.Close()

	kKafka := kafka.Writer{
		Addr:     kafka.TCP(""),
		Topic:    "price",
		Balancer: &kafka.LeastBytes{},
	}
	defer kKafka.Close()

	// БД

	//repoOutbox := repo.NewOutboxRepo(db)

	repoProd := repo.NewPostgresProductRepo(db)
	repoUser := repo.NewPostgresRepository(db)
	repoSubs := repo.NewPostgresSubscriptionRepository(db)

	// Сервисы

	serviceProd := sP.NewService(repoProd)
	serviceUser := sU.NewService(repoUser)
	serviceSubs := sS.NewService(repoSubs, rDB)

	// Хендлер

	handlerG := h.NewHandler(serviceSubs, serviceUser, serviceProd)

	// Parser Service

	parserConn, err := grpc.Dial(cfg.App.ParserAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к микросервису парсера: %v", err)
	}
	defer parserConn.Close()

	parserClient := par.NewGetPriceClient(parserConn)

	// Worker

	wetchPricer := w.NewPriceWatcher(repoProd, serviceSubs, l, parserClient, cb)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Server

	healthHandler := health.NewHealthCheck(rDB, db, l)

	mux := http.NewServeMux()

	healthHandler.RegisterRoutes(mux)

	healthServer := &http.Server{
		Addr:    cfg.Health.Port,
		Handler: mux,
	}

	go func() {
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("Ошибка HTTP сервера", slog.Any("error", err))
		}
	}()

	grpcServer := grpc.NewServer()

	g.RegisterPriceServiceServer(grpcServer, handlerG)
	reflection.Register(grpcServer)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("запуск воркера")
		if err := wetchPricer.Worker(ctx); err != nil {
			log.Println(err, "воркер завершился с ошибкой")
		}
	}()

	listener, err := net.Listen("tcp", cfg.App.GRPCServerPort)
	if err != nil {
		log.Fatalf("Error listening on port %s: %v", err, cfg.DB.Port)
	}

	go func() {
		log.Println("Запускаем gRPC сервер на порту :50051...")
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("Сервер остановлен: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("получен сигнал на прекрощение работы")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = healthServer.Shutdown(shutdownCtx)
	grpcServer.GracefulStop()

	log.Println("приложение закончило свою работу")
	wg.Wait()
}
