package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	FullName string `gorm:"not null" json:"full_name"`
	Password string `gorm:"not null" json:"-"` // Hidden in JSON responses
	Weight   int    `gorm:"not null" json:"weight"`
	Height   int    `gorm:"not null" json:"height"`

	Workouts     []Workout     `gorm:"foreignKey:UserID" json:"-"`
	ExerciseLogs []ExerciseLog `gorm:"foreignKey:UserID" json:"-"`
}
