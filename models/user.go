
package models

import "time"

type User struct {
	ID              uint      `gorm:"primaryKey"`
	Nombre          string    `gorm:"not null"`
	ApellidoPaterno string    `gorm:"not null"`
	ApellidoMaterno string    `gorm:"not null"`
	Correo          string    `gorm:"not null;unique"`
	Contraseña      string    `gorm:"not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}
