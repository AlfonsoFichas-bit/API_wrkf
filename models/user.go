package models

import "time"

// User represents a user account in the platform.
type User struct {
	ID              uint   `gorm:"primaryKey"`
	Nombre          string `gorm:"not null"`
	ApellidoPaterno string `gorm:"not null"`
	ApellidoMaterno string `gorm:"not null"`
	Correo          string `gorm:"not null;unique"`
	Contraseña      string `gorm:"not null"`
	// Role defines the user's platform-level role (e.g., 'user' or 'admin').
	Role      string    `gorm:"not null;default:'user'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
