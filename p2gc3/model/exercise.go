package model

type Exercise struct {
	ID          uint   `gorm:"primaryKey"`
	WorkoutID   uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`

	Workout      Workout       `gorm:"foreignKey:WorkoutID"`
	ExerciseLogs []ExerciseLog `gorm:"foreignKey:ExerciseID"`
}
