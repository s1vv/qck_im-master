package main

import (
	"qckim-backend/config"
	"qckim-backend/internal/handlers"
	"qckim-backend/internal/logger"
	"qckim-backend/internal/middleware"
	"qckim-backend/internal/repository"
	"qckim-backend/utils/email"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadConfig()

	cfg := config.GetConfig()

	logger.InitLogger()

	email.NewEmailSender()
	email.StartEmailWorker()

	db, err := repository.Connect(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPass)
	if err != nil {
		logger.Fatal("Main repo", err)
	}
	defer db.Close()

	// Выполнение миграций
	if err := repository.RunMigrations(db); err != nil {
		logger.Fatal("Main repo", err)

	}

	repo := repository.NewQckRepo(db)
	defer repo.Close()

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(repo)
	qckLinkHandler := handlers.NewQckLinkHandler(repo)

	// Создание Gin маршрутизатора
	r := gin.Default()

	// Настройка CORS
	r.Use(middleware.CORSConfig(cfg.AppEnv, cfg.BaseURL))
	r.Use(middleware.JWTAuth())

	// Регистрация маршрутов
	userHandler.RegisterRoutes(r)
	qckLinkHandler.RegisterRoutes(r)

	// Настройка порта из конфигурации или по умолчанию
	port := ":" + cfg.ServerPort

	// Запуск сервера
	logger.Info("Starting server on %s\n", "port", port)

	if err := r.Run(port); err != nil {
		logger.Fatal("Ошибка запуска сервера: %v", err)
	}

}
