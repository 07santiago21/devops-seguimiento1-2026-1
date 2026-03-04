package enrollment

import (
	"time"

	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/course"
	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/student"
)

type Enrollment struct {
	ID          string           `json:"id" gorm:"type:char(36);not null;primary_key"`
	StudentID   string           `json:"student_id" gorm:"type:char(36);not null"`
	Student     *student.Student `json:"student,omitempty" gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CourseID    string           `json:"course_id" gorm:"type:char(36);not null"`
	Course      *course.Course   `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	TotalAmount float64          `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time        `json:"created_at"`
}
