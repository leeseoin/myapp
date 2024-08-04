package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"myapp/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type GeminiService struct {
	Client *genai.Client
}

func NewGeminiService() (*GeminiService, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set in environment variables")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &GeminiService{Client: client}, nil
}

func (gs *GeminiService) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract text from form data
	requestText := r.FormValue("requestText")

	// Extract image files
	file1, file1Header, err := r.FormFile("image1")
	if err != nil {
		http.Error(w, "Unable to get image1", http.StatusBadRequest)
		return
	}
	defer file1.Close()

	file2, file2Header, err := r.FormFile("image2")
	if err != nil {
		http.Error(w, "Unable to get image2", http.StatusBadRequest)
		return
	}
	defer file2.Close()

	// Read files into byte slices
	imgData1, err := io.ReadAll(file1)
	if err != nil {
		http.Error(w, "Error reading image1", http.StatusInternalServerError)
		return
	}

	imgData2, err := io.ReadAll(file2)
	if err != nil {
		http.Error(w, "Error reading image2", http.StatusInternalServerError)
		return
	}

	// Detect image formats
	img1Format := detectImageFormat(file1Header.Filename)
	img2Format := detectImageFormat(file2Header.Filename)

	if img1Format == "" || img2Format == "" {
		http.Error(w, "Unsupported image format", http.StatusBadRequest)
		return
	}

	// Generate response from Gemini
	ctx := context.Background()
	model := gs.Client.GenerativeModel("gemini-1.5-flash")

	prompt := []genai.Part{
		genai.Text(requestText),
		genai.ImageData(img1Format, imgData1),
		genai.ImageData(img2Format, imgData2),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		http.Error(w, "Error generating content", http.StatusInternalServerError)
		return
	}

	// Send response back to client
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error marshalling response to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

// Detect image format based on file extension
func detectImageFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	default:
		return ""
	}
}

func (gs *GeminiService) AnalyzeNoteContentAndImage(ctx context.Context, note *model.Note, requestText string) (string, error) {
	// GenerativeModel 호출
	model := gs.Client.GenerativeModel("gemini-1.5-flash")

	prompt := []genai.Part{
		genai.Text(fmt.Sprintf("Analyze the content and image of the following note. Request: %s\nContent: %s\nImage URL: %s", requestText, note.Content, note.Img)),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", err
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(respJSON), nil
}
