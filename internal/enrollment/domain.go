package enrollment

import (
	"time"
)

type Enrollment struct {
	ID          string    `json:"id" gorm:"type:char(36);not null;primary_key"`
	StudentID   string    `json:"student_id" gorm:"type:char(36);not null"`
	CourseID    string    `json:"course_id" gorm:"type:char(36);not null"`
	TotalAmount float64   `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time `json:"created_at"`
}
