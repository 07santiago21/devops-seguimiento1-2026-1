package main

import (
	"log"
	"net/http"

	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/database"
	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/student"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {

	_ = godotenv.Load()

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := autoMigrate(db); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	repo := student.NewRepository(db)
	service := student.NewService(repo)
	handle := student.Handler(service)

	router := mux.NewRouter()
	router.HandleFunc("/students", handle.Create).Methods("POST")

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}

func autoMigrate(db *gorm.DB) error {
	return db.Session(&gorm.Session{PrepareStmt: false}).AutoMigrate(&student.Student{})
}
