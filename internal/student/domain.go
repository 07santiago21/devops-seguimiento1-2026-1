package student

import "time"

type Student struct {
	ID        string     `json:"id" gorm:"type:char(36);not null;primary_key;unique_index "`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	LastName  string     `json:"last_name" gorm:"type:varchar(100);not null"`
	Age       int32      `json:"age" gorm:"not null"`
	CreatedAt *time.Time `json:"-"`
}
