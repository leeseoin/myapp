package repository

import (
	"database/sql"
	"myapp/model"
)

// NoteRepository 구조체 정의
type NoteRepository struct {
	DB *sql.DB
}

// NewNoteRepository 함수 정의
func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{DB: db}
}

// Create 함수 정의
func (r *NoteRepository) Create(note *model.Note) (int, error) {
	result, err := r.DB.Exec("INSERT INTO notes (img, title, content, created_time, updated_time) VALUES (?, ?, ?, ?, ?)", note.Img, note.Title, note.Content, note.CreatedTime, note.UpdatedTime)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// GetByID 함수 정의
func (r *NoteRepository) GetByID(id int) (*model.Note, error) {
	row := r.DB.QueryRow("SELECT id, img, title, content, created_time, updated_time FROM notes WHERE id = ?", id)
	note := &model.Note{}
	err := row.Scan(&note.ID, &note.Img, &note.Title, &note.Content, &note.CreatedTime, &note.UpdatedTime)
	if err != nil {
		return nil, err
	}
	return note, nil
}

// GetAll 함수 정의
func (r *NoteRepository) GetAll() ([]*model.Note, error) {
	rows, err := r.DB.Query("SELECT id, img, title, content, created_time, updated_time FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*model.Note
	for rows.Next() {
		note := &model.Note{}
		err := rows.Scan(&note.ID, &note.Img, &note.Title, &note.Content, &note.CreatedTime, &note.UpdatedTime)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// Update 함수 정의
func (r *NoteRepository) Update(note *model.Note) error {
	_, err := r.DB.Exec("UPDATE notes SET img = ?, title = ?, content = ?, updated_time = ? WHERE id = ?",
		note.Img, note.Title, note.Content, note.UpdatedTime, note.ID)
	return err
}

// Delete 함수 정의
func (r *NoteRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}
