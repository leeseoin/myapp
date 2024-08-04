package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// GeminiResponse는 응답의 구조를 정의합니다.
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []string `json:"Parts"`
		} `json:"Content"`
	} `json:"Candidates"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access the API key from environment variables
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set in environment variables")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	// Initialize the chat
	// cs := model.StartChat()
	// cs.History = []*genai.Content{
	// 	{
	// 		Parts: []genai.Part{
	// 			genai.Text("Hello, I have 2 dogs in my house."),
	// 		},
	// 		Role: "user",
	// 	},
	// 	{
	// 		Parts: []genai.Part{
	// 			genai.Text("Great to meet you. What would you like to know?"),
	// 		},
	// 		Role: "model",
	// 	},
	// }

	// resp, err := cs.SendMessage(ctx, genai.Text("How many paws are in my house?"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// 이미지랑 텍스트 처리
	// 수정된 부분: 경로를 문자열 리터럴로 올바르게 작성
	imgData1, err := os.ReadFile("/Users/seoin/Desktop/note/myapp/uploads/real/kill.png")
	if err != nil {
		log.Fatal(err)
	}

	// imgData2, err := os.ReadFile("/Users/seoin/Desktop/note/myapp/uploads/real/dog.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := model.GenerateContent(
		ctx,
		genai.ImageData("png", imgData1),
		// genai.ImageData("png", imgData2), // 수정된 부분: 쉼표 추가
		genai.Text("do you know him?, if you know, tell his name, age, and where he came"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the response to JSON format for easier processing
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Error marshalling response to JSON:", err)
	}

	// Define a variable to hold the parsed response
	var response GeminiResponse
	if err := json.Unmarshal(respJSON, &response); err != nil {
		log.Fatal("Error unmarshalling response JSON:", err)
	}

	// Extract and print the content from the first candidate
	if len(response.Candidates) > 0 {
		if len(response.Candidates[0].Content.Parts) > 0 {
			// Print the first part of the content from the first candidate
			log.Println("Generated content:", response.Candidates[0].Content.Parts[0])
		} else {
			log.Println("No content parts found.")
		}
	} else {
		log.Println("No candidates found.")
	}
}
