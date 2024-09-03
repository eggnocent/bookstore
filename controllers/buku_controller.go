package controllers

import (
	"apimandiri/middlewares"
	"apimandiri/models"
	"apimandiri/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BukuController interface {
	AddUserBook(ctx *gin.Context)
	GetAllBooks(ctx *gin.Context)
	GetUserBook(ctx *gin.Context)
	GetBookByID(ctx *gin.Context)
	UpdateUserBook(ctx *gin.Context)
	UpdateBookByID(ctx *gin.Context)
	DeleteUserBook(ctx *gin.Context)
	DeleteBookByID(ctx *gin.Context)
	AddAuthorToBook(ctx *gin.Context)
	UpdateAuthorsForBook(ctx *gin.Context)
	DeleteAuthorFromBook(ctx *gin.Context)
}

type bukuController struct {
	service services.BookService
}

func NewBukuController(service services.BookService) BukuController {
	return &bukuController{service}
}

func (c *bukuController) AddUserBook(ctx *gin.Context) {
	id := ctx.Param("id")
	var input struct {
		NamaBuku  string `json:"NamaBuku"`
		TglTerbit string `json:"TglTerbit"`
		PenulisID uint   `json:"PenulisID"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		middlewares.HandleInternalServerError(ctx, "Input tidak valid")
		return
	}

	buku := models.Buku{
		NamaBuku:  input.NamaBuku,
		TglTerbit: input.TglTerbit,
		IdPenulis: input.PenulisID,
		UserID:    0,
	}

	// Panggil service untuk menambahkan buku
	err := c.service.AddBookToUser(id, buku)
	if err != nil {
		// Jika error yang terjadi adalah "user telah memiliki buku", tangani dengan spesifik
		if err.Error() == "user telah memiliki buku" {
			middlewares.HandleInternalServerError(ctx, "User telah memiliki buku, tidak dapat menambahkan lagi")
			return
		}

		// Tangani error lainnya dengan pesan yang lebih spesifik
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book added to user"})
}

// buku_controller.go
func (c *bukuController) GetAllBooks(ctx *gin.Context) {
	authorIDStr := ctx.Query("Penulis")
	userIDStr := ctx.Query("User")

	var books []models.Buku
	var err error
	var authorID uint
	var userID uint

	// Konversi ID penulis jika tersedia
	if authorIDStr != "" {
		authorIDParsed, err := strconv.ParseUint(authorIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID penulis tidak valid"})
			return
		}
		authorID = uint(authorIDParsed)
	}

	// Konversi ID user jika tersedia
	if userIDStr != "" {
		userIDParsed, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
			return
		}
		userID = uint(userIDParsed)
	}

	// Panggil service dengan parameter yang diperlukan
	books, err = c.service.GetAllBooks(authorID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan buku"})
		return
	}

	// Buat respons JSON untuk buku
	var booksResponse []gin.H
	for _, book := range books {
		var penulisMany []gin.H
		for _, penulis := range book.PenulisMany {
			penulisMany = append(penulisMany, gin.H{
				"ID":           penulis.ID,
				"NamaPenulis":  penulis.NamaPenulis,
				"EmailPenulis": penulis.EmailPenulis,
			})
		}

		bookResponse := gin.H{
			"ID":          book.ID,
			"NamaBuku":    book.NamaBuku,
			"TglTerbit":   book.TglTerbit,
			"PenulisMany": penulisMany,
			"User":        book.User.Username,
			"UserID":      book.UserID,
		}
		booksResponse = append(booksResponse, bookResponse)
	}

	ctx.JSON(http.StatusOK, booksResponse)
}

func (c *bukuController) GetBookByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID buku tidak valid"})
		return
	}

	book, err := c.service.GetBookByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	var penulisMany []gin.H
	for _, penulis := range book.PenulisMany {
		penulisMany = append(penulisMany, gin.H{
			"ID":           penulis.ID,
			"NamaPenulis":  penulis.NamaPenulis,
			"EmailPenulis": penulis.EmailPenulis,
		})
	}

	response := gin.H{
		"ID":          book.ID,
		"NamaBuku":    book.NamaBuku,
		"TglTerbit":   book.TglTerbit,
		"PenulisMany": penulisMany,
		"User":        book.User.Username,
		"UserID":      book.UserID,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *bukuController) GetUserBook(ctx *gin.Context) {
	id := ctx.Param("id") // Mengambil ID dari URL parameter
	buku, err := c.service.GetBookByUserID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ada"})
		return
	}

	response := map[string]interface{}{
		"ID":        buku.ID,
		"NamaBuku":  buku.NamaBuku,
		"Penulis":   buku.Penulis.NamaPenulis,
		"TglTerbit": buku.TglTerbit,
		"UserID":    buku.UserID,
		"User":      buku.User.Username,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *bukuController) UpdateUserBook(ctx *gin.Context) {
	id := ctx.Param("id")
	var buku models.Buku
	if err := ctx.ShouldBindJSON(&buku); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.service.UpdateBook(id, buku)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

func (c *bukuController) UpdateBookByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID buku tidak valid"})
		return
	}

	var book models.Buku
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = uint(id)

	err = c.service.UpdateBookByID(book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Buku berhasil diupdate"})
}

func (c *bukuController) DeleteUserBook(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteBook(id)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (c *bukuController) DeleteBookByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteBookByID(id)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pesan": "Buku telah terhapus"})
}

func (c *bukuController) AddAuthorToBook(ctx *gin.Context) {
	bookID := ctx.Param("book_id")
	authorID := ctx.Param("author_id")

	err := c.service.AddAuthorToBook(bookID, authorID)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pesan": "Author added to book successfully"})
}

func (c *bukuController) UpdateAuthorsForBook(ctx *gin.Context) {
	bookID := ctx.Param("book_id")
	var authorIDs []string
	if err := ctx.ShouldBindJSON(&authorIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.service.UpdateAuthorsForBook(bookID, authorIDs)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pesan": "Author Update Success for the book"})
}

func (c *bukuController) DeleteAuthorFromBook(ctx *gin.Context) {
	bookID := ctx.Param("book_id")
	authorID := ctx.Param("author_id")

	err := c.service.DeleteAuthorFromBook(bookID, authorID)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pesan": "author removed from book"})
}
