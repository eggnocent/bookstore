package repositories

import (
	"apimandiri/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetAllBooks(authorID uint, userID uint) ([]models.Buku, error)
	GetBookByUserID(userID uint) (*models.Buku, error)
	GetBookByID(id uint) (*models.Buku, error)
	AddBookToUser(book models.Buku) error
	UpdateBook(book models.Buku) error
	UpdateBookByID(book models.Buku) error
	DeleteBook(userID uint) error
	DeleteBookByID(id uint) error
	AddAuthorToBook(bookID uint, authorID uint) error
	UpdateAuthorsForBook(bookID uint, authors []models.Penulis) error
	DeleteAuthorFromBook(bookID uint, authorID uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (r *bookRepository) GetAllBooks(authorID uint, userID uint) ([]models.Buku, error) {
	var books []models.Buku
	query := r.db.Debug().Preload("PenulisMany").Preload("User")

	// Terapkan filter berdasarkan authorID jika disediakan
	if authorID != 0 {
		query = query.Where("id_penulis = ?", authorID)
	}

	// Terapkan filter berdasarkan userID jika disediakan
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	// Eksekusi query
	err := query.Find(&books).Error
	return books, err
}

func (r *bookRepository) GetBookByID(id uint) (*models.Buku, error) {
	var book models.Buku
	err := r.db.Debug().Preload("PenulisMany").Preload("User").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetBookByUserID(userID uint) (*models.Buku, error) {
	var book models.Buku
	err := r.db.Preload("User").Preload("Penulis").Where("user_id = ?", userID).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) AddBookToUser(book models.Buku) error {
	return r.db.Create(&book).Error
}

func (r *bookRepository) UpdateBook(book models.Buku) error {
	return r.db.Save(&book).Error
}

func (r *bookRepository) UpdateBookByID(book models.Buku) error {
	return r.db.Save(&book).Error
}

func (r *bookRepository) DeleteBook(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Buku{}).Error
}

func (r *bookRepository) DeleteBookByID(id uint) error {
	return r.db.Delete(&models.Buku{}, id).Error
}

func (r *bookRepository) AddAuthorToBook(bookID uint, authorID uint) error {
	var book models.Buku
	if err := r.db.Preload("PenulisMany").First(&book, bookID).Error; err != nil {
		return err
	}

	var author models.Penulis
	if err := r.db.First(&author, authorID).Error; err != nil {
		return err
	}

	return r.db.Model(&book).Association("PenulisMany").Append(&author)
}

func (r *bookRepository) UpdateAuthorsForBook(bookID uint, authors []models.Penulis) error {
	var book models.Buku
	if err := r.db.Preload("PenulisMany").First(&book, bookID).Error; err != nil {
		return err
	}

	return r.db.Model(&book).Association("PenulisMany").Replace(authors)
}

func (r *bookRepository) DeleteAuthorFromBook(bookID, authorID uint) error {
	var book models.Buku
	if err := r.db.Preload("PenulisMany").First(&book, bookID).Error; err != nil {
		return err
	}

	var author models.Penulis
	if err := r.db.First(&author, authorID).Error; err != nil {
		return err
	}
	return r.db.Model(&book).Association("PenulisMany").Delete(&author)
}
