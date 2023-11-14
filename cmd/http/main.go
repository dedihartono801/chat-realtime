package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/dedihartono801/chat-realtime/cmd/http/routes"
	"github.com/dedihartono801/chat-realtime/database"
	"github.com/dedihartono801/chat-realtime/internal/app/queue/kafka"
	"github.com/dedihartono801/chat-realtime/internal/app/repository"
	"github.com/dedihartono801/chat-realtime/internal/app/usecase/chat"
	"github.com/dedihartono801/chat-realtime/internal/app/usecase/user"
	"github.com/dedihartono801/chat-realtime/internal/delivery/http"
	"github.com/dedihartono801/chat-realtime/pkg/config"
	"github.com/dedihartono801/chat-realtime/pkg/validator"
	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var db *sql.DB

	// Number of retry attempts
	maxRetries := 5
	retryInterval := 5 * time.Second

	for retries := 1; retries <= maxRetries; retries++ {
		fmt.Printf("Attempt %d to connect to the database...\n", retries)

		// Attempt to connect to the database
		db, err = database.InitPostgres()

		if err == nil {
			// Connection successful, break out of the loop
			fmt.Println("Connected to the database!")
			break
		}

		// Connection failed, wait for a short interval before retrying
		fmt.Printf("Error connecting to the database: %v\n", err)
		fmt.Printf("Retrying in %s...\n", retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		// All retry attempts failed
		log.Fatalf("Failed to connect to the database after %d attempts. Error: %v\n", maxRetries, err)
	}

	kafkaProducer, err := kafka.NewKafkaProducer(config.GetEnv("KAFKA_ADDRESS"), config.GetEnv("CHAT_TOPIC"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	validator := validator.NewValidator(validatorv10.New())
	userRepository := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepository, validator)
	userHandler := http.NewUserHandler(userService)

	userChatRepository := repository.NewUserChatRepository(db)
	messageRepository := repository.NewMessageRepository(db)
	chatService := chat.NewChatService(userChatRepository, messageRepository, kafkaProducer, validator)
	chatHandler := http.NewChatHandler(chatService)

	app := fiber.New()
	routes.SetupRoutes(app)
	routes.UserRouter(app, userHandler)
	routes.ChatRouter(app, chatHandler)

	if err := app.Listen(":5001"); err != nil {
		log.Fatalf("listen: %s", err)
	}

}
