package models

type BukuPenulis struct {
	BukuID    uint `gorm:"primaryKey"`
	PenulisID uint `gorm:"primaryKey"`
}
