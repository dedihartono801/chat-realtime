package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dedihartono801/chat-realtime/pkg/config"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
)

func GetMessageFromSupabase(query string) ([]*dto.FetchMessageDto, error) {
	var data []*dto.FetchMessageDto

	// Create a new HTTP POST request
	req, err := http.NewRequest("GET", config.GetEnv("SUPABASE_PROJECT_URL")+"/rest/v1/"+config.GetEnv("SUPABASE_TABLE")+"?select=id,user_chat_id,message_text,created_at&or=("+query+")", nil)
	if err != nil {
		return nil, err
	}

	// Set headers for the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", config.GetEnv("SUPABASE_KEY"))
	req.Header.Set("Authorization", "Bearer "+config.GetEnv("SUPABASE_KEY")) // Set your specific authorization token
	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response
	// Response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	fmt.Println("Data:", string(body))

	return data, nil

}
