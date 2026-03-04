package course

import "time"

type Course struct {
	ID          string     `json:"id" gorm:"type:char(36);not null;primary_key;unique_index"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null"`
	Description string     `json:"description" gorm:"type:text"`
	Credits     int32      `json:"credits" gorm:"not null"`
	Capacity    int32      `json:"capacity" gorm:"not null"`
	CreatedAt   *time.Time `json:"-"`
}
