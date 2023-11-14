package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dedihartono801/chat-realtime/database"
	"github.com/dedihartono801/chat-realtime/internal/app/repository"
	"github.com/dedihartono801/chat-realtime/internal/app/usecase/supabase"
	"github.com/dedihartono801/chat-realtime/internal/delivery/websocket"
	"github.com/dedihartono801/chat-realtime/pkg/config"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
	"github.com/dedihartono801/chat-realtime/pkg/validator"
	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	realtimego "github.com/overseedio/realtime-go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var postgresDB *sql.DB

	// Number of retry attempts
	maxRetries := 5
	retryInterval := 5 * time.Second

	for retries := 1; retries <= maxRetries; retries++ {
		fmt.Printf("Attempt %d to connect to the database...\n", retries)

		// Attempt to connect to the database
		postgresDB, err = database.InitPostgres()

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

	validator := validator.NewValidator(validatorv10.New())
	messageRepository := repository.NewMessageRepository(postgresDB)
	supabaseService := supabase.NewSupabaseService(messageRepository, validator)
	supabaseHandler := websocket.NewSupabase(supabaseService)

	// (optional) auth token
	//const RLS_TOKEN = ""

	// create client
	c, err := realtimego.NewClient(config.GetEnv("SUPABASE_PROJECT_URL"), config.GetEnv("SUPABASE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// connect to server
	err = c.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// create and subscribe to channel
	db := "realtime"
	schema := "public"
	table := config.GetEnv("SUPABASE_TABLE")
	ch, err := c.Channel(realtimego.WithTable(&db, &schema, &table))
	if err != nil {
		log.Fatal(err)
	}

	//dt := make(map[string]string)
	var data dto.InsertInfo

	// setup hooks
	ch.OnInsert = func(m realtimego.Message) {
		// Assuming data holds the map data
		jsonData, err := json.Marshal(m.Payload)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			log.Fatal(err)
		}

		req := &dto.SaveMessageDto{
			UserChatID:  data.Record.UserChatID,
			MessageText: data.Record.MessageText,
		}

		err = supabaseHandler.SaveMessages(req)
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	// ch.OnDelete = func(m realtimego.Message) {
	// 	log.Println("***ON DELETE....", m)
	// }
	// ch.OnUpdate = func(m realtimego.Message) {
	// 	log.Println("***ON UPDATE....", m)
	// }

	// subscribe to channel
	err = ch.Subscribe()
	if err != nil {
		log.Fatal(err)
	}

	// This will keep the main function running
	select {}

}
