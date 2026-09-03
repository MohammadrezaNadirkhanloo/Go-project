package handlers

import (
	"net/http"
	"strconv"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/dto"
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/helper"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/services"
	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	services *services.CountryService
}

func NewCountryHandler(cfg *config.Config) *CountryHandler {
	return &CountryHandler{
		services: services.NewCountryService(cfg),
	}
}

// Create
func (h *CountryHandler) Create(c *gin.Context) {
	req := dto.CreateAndUpdateCountryRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	res, err := h.services.Create(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, 0))
}

// UpdateCountry godoc
// @Summary Update a country
// @Description Update a country
// @Tags Countries
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Param Request body dto.CreateUpdateCountryRequest true "Update a country"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.CountryResponse} "Country response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/countries/{id} [put]
// @Security AuthBearer
// Update
func (h *CountryHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	req := dto.CreateAndUpdateCountryRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	res, err := h.services.Update(c, id, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, 0))
}

// DeleteCountry godoc
// @Summary Delete a country
// @Description Delete a country
// @Tags Countries
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} helper.BaseHttpResponse "response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/countries/{id} [delete]
// @Security AuthBearer
// Delete
func (h *CountryHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, helper.GenerateBaseResponse(nil, false, -1))
		return
	}
	err := h.services.Delete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, 0))
}

// GetCountry godoc
// @Summary Get a country
// @Description Get a country
// @Tags Countries
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.CountryResponse} "Country response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/countries/{id} [get]
// @Security AuthBearer
// Get
func (h *CountryHandler) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, helper.GenerateBaseResponse(nil, false, -1))
		return
	}
	res, err := h.services.GetById(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, 0))
}

// GetCountries godoc
// @Summary Get Countries
// @Description Get Countries
// @Tags Countries
// @Accept json
// @produces json
// @Param Request body dto.PaginationInputWithFilter true "Request"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.PagedList[dto.CountryResponse]} "Country response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/countries/get-by-filter [post]
// @Security AuthBearer
func (h *CountryHandler) GetByFilter(c *gin.Context) {
	req := dto.PaginationInputWithFilter{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	res, err := h.services.GetByFilter(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, 0))
}
