package service

import (
	"myapp/model"
	"myapp/repository"
	"time"
)

// NoteService 구조체 정의
type NoteService struct {
	Repo *repository.NoteRepository
}

// NewNoteService 함수 정의
func NewNoteService(repo *repository.NoteRepository) *NoteService {
	return &NoteService{Repo: repo}
}

// CreateNote 함수 정의
func (s *NoteService) CreateNote(title, content, img string) (*model.Note, error) {
	now := time.Now()
	note := &model.Note{
		Title:       title,
		Content:     content,
		Img:         img,
		CreatedTime: now,
		UpdatedTime: nil,
	}
	id, err := s.Repo.Create(note)
	if err != nil {
		return nil, err
	}
	note.ID = id
	return note, nil
}

// GetAllNotes 함수 정의
func (s *NoteService) GetAllNotes() ([]*model.Note, error) {
	return s.Repo.GetAll()
}

// GetNoteByID 함수 정의
func (s *NoteService) GetNoteByID(id int) (*model.Note, error) {
	return s.Repo.GetByID(id)
}

// UpdateNote 함수 정의
func (s *NoteService) UpdateNote(id int, title, content, img string) (*model.Note, error) {
	now := time.Now()
	note := &model.Note{
		ID:          id,
		Title:       title,
		Content:     content,
		Img:         img,
		UpdatedTime: &now,
	}

	err := s.Repo.Update(note)
	if err != nil {
		return nil, err
	}

	// 업데이트된 노트를 다시 조회
	updatedNote, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updatedNote, nil
}

// DeleteNote 함수 정의
func (s *NoteService) DeleteNote(id int) error {
	return s.Repo.Delete(id)
}
