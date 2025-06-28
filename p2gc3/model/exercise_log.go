package model

import "time"

type ExerciseLog struct {
	ID         uint      `gorm:"primaryKey"`
	ExerciseID uint      `gorm:"not null"` // FK to Exercise
	UserID     uint      `gorm:"not null"` // FK to User
	SetCount   int       `gorm:"not null"`
	RepCount   int       `gorm:"not null"`
	Weight     int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`

	User     User     `gorm:"foreignKey:UserID"`
	Exercise Exercise `gorm:"foreignKey:ExerciseID"`
}
