package controllers

import (
	"apimandiri/middlewares"
	"apimandiri/models"
	"apimandiri/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PenulisController interface {
	CreatePenulis(ctx *gin.Context)
	GetAllPenulis(ctx *gin.Context)
	GetPenulisByID(ctx *gin.Context)
	UpdatePenulis(ctx *gin.Context)
	DeletePenulis(ctx *gin.Context)
}

type penulisControllerImpl struct {
	service services.PenulisService
}

func NewPenulisController(service services.PenulisService) PenulisController {
	return &penulisControllerImpl{service}
}

type BukuResponse struct {
	NamaBuku  string `json:"NamaBuku"`
	TglTerbit string `json:"TglTerbit"`
}

type PenulisResponse struct {
	ID           uint           `json:"ID"`
	NamaPenulis  string         `json:"NamaPenulis"`
	EmailPenulis string         `json:"EmailPenulis"`
	Buku         []BukuResponse `json:"Buku"`
}

func (c *penulisControllerImpl) CreatePenulis(ctx *gin.Context) {
	var penulis models.Penulis
	if err := ctx.ShouldBindJSON(&penulis); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreatePenulis(penulis); err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"pesan": "Penulis berhasil dibuat"})
}

func (c *penulisControllerImpl) GetAllPenulis(ctx *gin.Context) {
	penulisList, err := c.service.GetAllPenulis()
	if err != nil {
		middlewares.HandleInternalServerError(ctx, "failed to fetch author")
		return
	}

	var penulisResponses []PenulisResponse

	for _, penulis := range penulisList {
		var bukuResponses []BukuResponse
		for _, buku := range penulis.Buku {
			bukuResponses = append(bukuResponses, BukuResponse{
				NamaBuku:  buku.NamaBuku,
				TglTerbit: buku.TglTerbit,
			})
		}

		penulisResponse := PenulisResponse{
			ID:           penulis.ID,
			NamaPenulis:  penulis.NamaPenulis,
			EmailPenulis: penulis.EmailPenulis,
			Buku:         bukuResponses,
		}

		penulisResponses = append(penulisResponses, penulisResponse)
	}

	ctx.JSON(http.StatusOK, penulisResponses)
}

func (c *penulisControllerImpl) GetPenulisByID(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	penulis, err := c.service.GetPenulisByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Penulis tidak ada"})
		return
	}
	var bukuResponses []BukuResponse
	for _, buku := range penulis.Buku {
		bukuResponses = append(bukuResponses, BukuResponse{
			NamaBuku:  buku.NamaBuku,
			TglTerbit: buku.TglTerbit,
		})
	}
	response := PenulisResponse{
		ID:           penulis.ID,
		NamaPenulis:  penulis.NamaPenulis,
		EmailPenulis: penulis.EmailPenulis,
		Buku:         bukuResponses,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *penulisControllerImpl) UpdatePenulis(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	var penulis models.Penulis
	if err := ctx.ShouldBindJSON(&penulis); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	penulis.ID = uint(id)
	if err := c.service.UpdatePenulis(penulis); err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pesan": "penulis berhasil diupdate"})
}

func (c *penulisControllerImpl) DeletePenulis(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeletePenulis(id)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pesan": "Penulis berhasil di hapus"})
}
