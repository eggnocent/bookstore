package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
	FullName string `gorm:"not null"`
	CreateAt time.Time
	UpdateAt time.Time
	Buku     *Buku `gorm:"foreignKey:UserID"`
}
