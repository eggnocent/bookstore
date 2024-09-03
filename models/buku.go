package models

type Buku struct {
	ID          uint      `gorm:"primaryKey"`
	NamaBuku    string    `gorm:"not null"`
	TglTerbit   string    `gorm:"not null"`
	IdPenulis   uint      `gorm:"foreignKey:IdPenulis;constraint:OnDelete:CASCADE;"`
	Penulis     Penulis   `gorm:"foreignKey:IdPenulis"`
	PenulisMany []Penulis `gorm:"many2many:buku_penulis;"`
	UserID      uint      `gorm:"unique"`
	User        *User     `gorm:"foreignKey:UserID"`
}
