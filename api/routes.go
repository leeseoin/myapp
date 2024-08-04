package api

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, noteHandler *NoteHandler) {
	e.POST("/notes", noteHandler.CreateNoteHandler)
	e.GET("/notes/:id", noteHandler.GetNoteByIDHandler)
	e.GET("/notes/all", noteHandler.GetAllNotesHandler)
	e.PUT("/notes/:id", noteHandler.UpdateNoteHandler)
	e.DELETE("/notes/:id", noteHandler.DeleteNoteHandler)
	e.POST("/api/notes/:id/analyze", noteHandler.AnalyzeNoteHandler)
}
