package api

import (
	"context"
	"log"
	"myapp/model"
	"myapp/service"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// NoteHandler 구조체 정의
type NoteHandler struct {
	NoteService   *service.NoteService
	GeminiService *service.GeminiService
}

// NewNoteHandler 함수 정의
func NewNoteHandler(noteService *service.NoteService, geminiService *service.GeminiService) *NoteHandler {
	return &NoteHandler{
		NoteService:   noteService,
		GeminiService: geminiService,
	}
}

// NoteResponse 구조체 정의
type NoteResponse struct {
	ID          int     `json:"id"`
	Img         string  `json:"img"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	CreatedTime string  `json:"created_time"`
	UpdatedTime *string `json:"updated_time"`
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func formatOptionalTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := formatTime(*t)
	return &s
}

type CreateNoteRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Img       string `json:"img"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateNoteHandler 함수 정의
func (h *NoteHandler) CreateNoteHandler(c echo.Context) error {
	// JSON 데이터 수신
	var req CreateNoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error message": "Invalid request format",
		})
	}

	// 기본 값 설정
	img := req.Img
	if img == "" {
		img = "" // 서버에서 빈 문자열로 설정
	}

	// 노트 생성
	note, err := h.NoteService.CreateNote(req.Title, req.Content, img)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error message": err.Error(),
		})
	}
	// 로그에 데이터 출력
	log.Printf("Received request: %+v\n", req)

	// 응답 생성
	response := map[string]interface{}{
		"message": "Note created successfully",
		"note_info": NoteResponse{
			ID:          note.ID,
			Img:         note.Img,
			Title:       note.Title,
			Content:     note.Content,
			CreatedTime: formatTime(note.CreatedTime),
			UpdatedTime: formatOptionalTime(note.UpdatedTime),
		},
	}

	return c.JSON(http.StatusCreated, response)
}

// UpdateNoteHandler 함수 정의
func (h *NoteHandler) UpdateNoteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error message": "Invalid ID format",
		})
	}

	// JSON 데이터 파싱
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Img     string `json:"img"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error message": "Invalid request format",
		})
	}

	// 노트 업데이트
	note, err := h.NoteService.UpdateNote(id, req.Title, req.Content, req.Img)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error message": err.Error(),
		})
	}

	// 응답 생성
	response := map[string]interface{}{
		"message": "Note updated successfully",
		"note_info": NoteResponse{
			ID:          note.ID,
			Img:         note.Img,
			Title:       note.Title,
			Content:     note.Content,
			CreatedTime: formatTime(note.CreatedTime),
			UpdatedTime: formatOptionalTime(note.UpdatedTime),
		},
	}

	return c.JSON(http.StatusOK, response)
}

// GetAllNotesHandler 함수 정의(노트 싹다 가져오기)
func (h *NoteHandler) GetAllNotesHandler(c echo.Context) error {
	notes, err := h.NoteService.GetAllNotes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error message": err.Error(),
		})
	}

	// 응답 생성
	response := map[string]interface{}{
		"message": "Notes retrieved successfully",
		"notes":   notesToResponse(notes),
	}

	return c.JSON(http.StatusOK, response)
}

// Note를 NoteResponse로 변환하는 함수
func notesToResponse(notes []*model.Note) []NoteResponse {
	responses := make([]NoteResponse, len(notes))
	for i, note := range notes {
		responses[i] = NoteResponse{
			ID:          note.ID,
			Img:         note.Img,
			Title:       note.Title,
			Content:     note.Content,
			CreatedTime: formatTime(note.CreatedTime),
			UpdatedTime: formatOptionalTime(note.UpdatedTime),
		}
	}
	return responses
}

// GetNoteByIDHandler 함수 정의(원하는 id에 해당하는 노트만 가져오기)
func (h *NoteHandler) GetNoteByIDHandler(c echo.Context) error {
	// URL에서 ID 추출
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error message": err.Error(),
		})
	}
	// 서비스에서 ID에 해당하는 노트 조회
	note, err := h.NoteService.GetNoteByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error message": err.Error(),
		})
	}

	// 응답 생성
	response := map[string]interface{}{
		"message": "Note retrieved successfully",
		"note_info": NoteResponse{
			ID:          note.ID,
			Img:         note.Img,
			Title:       note.Title,
			Content:     note.Content,
			CreatedTime: formatTime(note.CreatedTime),
			UpdatedTime: formatOptionalTime(note.UpdatedTime),
		},
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteNoteHandler 함수 정의(물리 삭제임)
func (h *NoteHandler) DeleteNoteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error message": "Invalid ID format",
		})
	}

	// 삭제할 노트 조회
	_, err = h.NoteService.GetNoteByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error message": err.Error(),
		})
	}

	// 서비스 레이어에서 노트 삭제
	err = h.NoteService.DeleteNote(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Note deleted successfully",
	})
}

// AnalyzeNoteHandler 함수 정의
func (h *NoteHandler) AnalyzeNoteHandler(c echo.Context) error {
	// ID 추출
	idParam := c.Param("id")

	// 노트 가져오기
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error message": "Invalid ID format",
		})
	}

	// 노트 데이터 가져오기
	note, err := h.NoteService.GetNoteByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error message": err.Error(),
		})
	}

	// 요청 데이터 추출
	requestText := c.FormValue("request")
	if requestText == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error message": "Request text is required",
		})
	}

	// Gemini API 호출
	analysisResult, err := h.GeminiService.AnalyzeNoteContentAndImage(context.Background(), note, requestText)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error message": err.Error(),
		})
	}

	// 응답 생성
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Note analyzed successfully",
		"result":  analysisResult,
	})
}
