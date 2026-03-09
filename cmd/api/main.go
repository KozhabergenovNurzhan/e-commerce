package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"ecommerce/internal/config"
	"ecommerce/internal/repository"
	"ecommerce/internal/server"
	"ecommerce/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	db, err := repository.NewPostgres(cfg.DSN())
	if err != nil {
		log.Error("failed to connect to postgres", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	log.Info("connected to postgres")

	// repos
	userRepo := repository.NewUserRepo(db)
	productRepo := repository.NewProductRepo(db)
	orderRepo := repository.NewOrderRepo(db)
	cartRepo := repository.NewCartRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)

	// services
	svc := service.NewServices(
		service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.TTL),
		service.NewProductService(productRepo),
		service.NewCartService(cartRepo, productRepo),
		service.NewOrderService(orderRepo, cartRepo, productRepo),
		service.NewCategoryService(categoryRepo),
	)

	r := server.New(svc, log, cfg.JWT.Secret)

	log.Info("starting server", slog.String("addr", cfg.HTTPAddr))

	go func() {
		if err := r.Run(cfg.HTTPAddr); err != nil {
			log.Error("server stopped", slog.String("err", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down...")
}
