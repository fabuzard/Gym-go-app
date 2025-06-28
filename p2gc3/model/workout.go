package model

type Workout struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	UserID      uint       `gorm:"not null"`
	User        User       `gorm:"foreignKey:UserID"`
	Exercises   []Exercise `gorm:"foreignKey:WorkoutID"`
}
