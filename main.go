package main

import (
	"log"
	"net/http"
	"os"

	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/course"
	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/database"
	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/enrollment"
	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/student"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {

	_ = godotenv.Load()

	log.Println("Connecting to DB...")
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	log.Println("Connected to DB")

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
	router.HandleFunc("/api/v2/students", handle.CreateV2).Methods("POST")

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
	router.HandleFunc("/enrollments/{id}", enrollEnd.Patch).Methods("PATCH")
	router.HandleFunc("/enrollments/{id}", enrollEnd.Put).Methods("PUT")

	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		log.Println("Running in Lambda")
		adapter := gorillamux.New(router)
		lambda.Start(adapter.ProxyWithContext)
	} else {
		log.Println("Running locally on :8000")
		if err := http.ListenAndServe(":8000", router); err != nil {
			log.Fatal(err)
		}
	}

}

func autoMigrate(db *gorm.DB) error {
	return db.Session(&gorm.Session{PrepareStmt: false}).AutoMigrate(&student.Student{}, &course.Course{}, enrollment.Enrollment{})
}
