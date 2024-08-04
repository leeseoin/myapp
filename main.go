package main

import (
	"database/sql"
	"log"
	"myapp/api"
	"myapp/config"
	"myapp/repository"
	"myapp/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Echo 웹 프레임워크 인스턴스 생성
	e := echo.New()

	// CORS 미들웨어 설정
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// 환경 변수 로드
	cfg := config.LoadConfig()

	// SQLite 데이터베이스 연결
	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	// 데이터베이스 스키마 초기화
	err = initializeSchema(db)
	if err != nil {
		log.Fatalf("could not initialize database schema: %v", err)
	}

	// 레포지토리, 서비스, 핸들러 생성
	repo := repository.NewNoteRepository(db)
	noteService := service.NewNoteService(repo)
	geminiService, err := service.NewGeminiService()
	if err != nil {
		log.Fatalf("could not initialize Gemini service: %v", err)
	}
	noteHandler := api.NewNoteHandler(noteService, geminiService)

	// 라우팅 설정
	api.RegisterRoutes(e, noteHandler)

	// 이미지 핸들러
	e.Static("/uploads", "myapp/uploads/.cache")

	// 서버 시작
	e.Logger.Fatal(e.Start(":8080"))
}

// 데이터베이스 스키마 초기화 함수
func initializeSchema(db *sql.DB) error {
	// 노트 테이블 생성 쿼리
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        content TEXT NOT NULL,
        img TEXT,
        created_time DATETIME NOT NULL,
        updated_time DATETIME
    );
    `

	// 테이블 생성 쿼리 실행
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}
