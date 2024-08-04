package model

import "time"

// Note 구조체 정의
type Note struct {
	ID          int        `json:"id"`
	Img         string     `json:"img"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	CreatedTime time.Time  `json:"created_time"`
	UpdatedTime *time.Time `json:"updated_time"`
}
