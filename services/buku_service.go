package services

import (
	"apimandiri/models"
	"apimandiri/repositories"
	"errors"
	"strconv"
	"strings"
)

type BookService interface {
	GetBookByUserID(userID string) (*models.Buku, error)
	GetBookByID(id uint) (*models.Buku, error)
	GetAllBooks(authorID uint, userID uint) ([]models.Buku, error)
	AddBookToUser(userID string, book models.Buku) error
	UpdateBook(userID string, book models.Buku) error
	UpdateBookByID(book models.Buku) error
	DeleteBook(userID string) error
	DeleteBookByID(id string) error
	AddAuthorToBook(bookID, authorID string) error
	UpdateAuthorsForBook(bookID string, authorIDs []string) error
	DeleteAuthorFromBook(bookID, authorID string) error
}

type bookServiceImpl struct {
	repo repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) BookService {
	return &bookServiceImpl{repo}
}

func (s *bookServiceImpl) GetAllBooks(authorID uint, userID uint) ([]models.Buku, error) {
	return s.repo.GetAllBooks(authorID, userID)
}

func (s *bookServiceImpl) GetBookByID(id uint) (*models.Buku, error) {
	return s.repo.GetBookByID(id)
}

func (s *bookServiceImpl) GetBookByUserID(userID string) (*models.Buku, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetBookByUserID(uint(id))
}

func (s *bookServiceImpl) AddBookToUser(userID string, book models.Buku) error {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("ID user tidak valid")
	}
	existingbook, err := s.repo.GetBookByUserID(uint(id))
	if err == nil && existingbook != nil {
		return errors.New("user telah memiliki buku")
	}
	book.UserID = uint(id)

	if book.IdPenulis == 0 {
		return errors.New("penulis ID harus valid")
	}

	return s.repo.AddBookToUser(book)
}

func (s *bookServiceImpl) UpdateBook(userID string, book models.Buku) error {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}

	existingBook, err := s.repo.GetBookByUserID(uint(id))
	if err != nil {
		return err
	}

	existingBook.NamaBuku = book.NamaBuku
	existingBook.IdPenulis = book.IdPenulis
	existingBook.TglTerbit = book.TglTerbit

	return s.repo.UpdateBook(*existingBook)
}

func (s *bookServiceImpl) UpdateBookByID(book models.Buku) error {
	existingBook, err := s.repo.GetBookByID(book.ID)
	if err != nil {
		return errors.New("buku tidak ditemukan")
	}

	existingBook.NamaBuku = book.NamaBuku
	existingBook.TglTerbit = book.TglTerbit
	existingBook.IdPenulis = book.IdPenulis

	return s.repo.UpdateBookByID(*existingBook)
}

func (s *bookServiceImpl) DeleteBook(userID string) error {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}
	return s.repo.DeleteBook(uint(id))
}

func (s *bookServiceImpl) DeleteBookByID(id string) error {
	bookID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("ID Buku tidak valid")
	}
	return s.repo.DeleteBookByID(uint(bookID))
}

func (s *bookServiceImpl) AddAuthorToBook(bookID, authorID string) error {
	bID, err := strconv.ParseUint(bookID, 10, 32)
	if err != nil {
		return errors.New("invalid book ID")
	}

	aID, err := strconv.ParseUint(authorID, 10, 32)
	if err != nil {
		return errors.New("invalid author ID")
	}

	return s.repo.AddAuthorToBook(uint(bID), uint(aID))
}

func (s *bookServiceImpl) UpdateAuthorsForBook(bookID string, authorIDs []string) error {
	bID, err := strconv.ParseUint(bookID, 10, 32)
	if err != nil {
		return errors.New("invalid book ID")
	}
	var penulisMany []models.Penulis
	for _, authorID := range authorIDs {
		aID, err := strconv.ParseUint(authorID, 10, 32)
		if err != nil {
			return errors.New("invalid author ID")
		}
		penulisMany = append(penulisMany, models.Penulis{ID: uint(aID)})
	}
	return s.repo.UpdateAuthorsForBook(uint(bID), penulisMany)
}

// services/buku_service.go
func (s *bookServiceImpl) DeleteAuthorFromBook(bookID, authorID string) error {
	bID, err := strconv.ParseUint(bookID, 10, 32)
	if err != nil {
		return errors.New("invalid book ID")
	}

	aID, err := strconv.ParseUint(authorID, 10, 32)
	if err != nil {
		return errors.New("invalid author ID")
	}

	err = s.repo.DeleteAuthorFromBook(uint(bID), uint(aID))
	if err != nil {
		// Deteksi pesan error foreign key dari MySQL
		if strings.Contains(err.Error(), "foreign key constraint fails") || strings.Contains(err.Error(), "1451") {
			return errors.New("Tidak dapat menghapus user, karena memiliki relasi dengan tabel lain.")
		}
		return err
	}
	return nil
}
