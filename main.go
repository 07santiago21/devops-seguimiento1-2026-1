package main

import (
	"log"
	"net/http"

	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/course"
	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/database"
	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/enrollment"
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

	// Course setup
	courseRepo := course.NewRepository(db)
	courseService := course.NewService(courseRepo)
	courseHandle := course.Handler(courseService)

	enrollRepo := enrollment.NewRepository(db)
	enrollSvc := enrollment.NewService(enrollRepo)
	enrollEnd := enrollment.MakeEndpoints(enrollSvc)

	router := mux.NewRouter()
	router.HandleFunc("/students", handle.Create).Methods("POST")
	router.HandleFunc("/students", handle.GetAll).Methods("GET")
	router.HandleFunc("/students/{id}", handle.Get).Methods("GET")
	router.HandleFunc("/students/{id}", handle.Delete).Methods("DELETE")
	router.HandleFunc("/students/{id}", handle.Patch).Methods("PATCH")
	router.HandleFunc("/students/{id}", handle.Put).Methods("PUT")

	router.HandleFunc("/courses", courseHandle.Create).Methods("POST")
	router.HandleFunc("/courses", courseHandle.GetAll).Methods("GET")
	router.HandleFunc("/courses/{id}", courseHandle.Get).Methods("GET")
	router.HandleFunc("/courses/{id}", courseHandle.Delete).Methods("DELETE")
	router.HandleFunc("/courses/{id}", courseHandle.Patch).Methods("PATCH")
	router.HandleFunc("/courses/{id}", courseHandle.Put).Methods("PUT")

	router.HandleFunc("/enrollments", enrollEnd.Create).Methods("POST")
	router.HandleFunc("/enrollments", enrollEnd.GetAll).Methods("GET")
	router.HandleFunc("/enrollments/{id}", enrollEnd.Get).Methods("GET")
	router.HandleFunc("/enrollments/{id}", enrollEnd.Delete).Methods("DELETE")
	router.HandleFunc("/enrollments/{id}", enrollEnd.GetAll).Methods("PATCHT")
	router.HandleFunc("/enrollments/{id}", enrollEnd.GetAll).Methods("PUT")

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}

func autoMigrate(db *gorm.DB) error {
	return db.Session(&gorm.Session{PrepareStmt: false}).AutoMigrate(&student.Student{}, &course.Course{}, enrollment.Enrollment{})
}
