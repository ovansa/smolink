package controller

import (
	"net/http"
	"smolink/internal/errors"
	"smolink/internal/service"

	"github.com/gin-gonic/gin"
)

type URLController struct {
	service *service.URLService
}

func NewURLController(service *service.URLService) *URLController {
	return &URLController{service: service}
}

func (uc *URLController) ShortenURL(c *gin.Context) {
	var payload struct {
		URL        string `json:"url"`
		CustomCode string `json:"customCode"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	result, err := uc.service.ShortenURL(c, payload.URL, payload.CustomCode)
	if err != nil {
		apiErr := errors.ExtractAPIError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"shortCode": result.ShortCode, "originalUrl": result.OriginalURL})
}

func (uc *URLController) ResolveURL(c *gin.Context) {
	code := c.Param("code")
	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	original, err := uc.service.ResolveURL(c, code, ip, ua)
	if err != nil {
		apiErr := errors.ExtractAPIError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.Redirect(http.StatusFound, original)
}
