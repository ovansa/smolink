package controller

import (
	"log"
	"net/http"
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
		CustomCode string `json:"custom_code"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	result, err := uc.service.ShortenURL(c, payload.URL, payload.CustomCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_code": result.ShortCode, "original_url": result.OriginalURL})
}

func (uc *URLController) ResolveURL(c *gin.Context) {
	code := c.Param("code")
	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	log.Printf("Resolving URL for code: %s, IP: %s, User-Agent: %s", code, ip, ua)

	original, err := uc.service.ResolveURL(c, code, ip, ua)
	log.Printf("Original URL: %s", original)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	c.Redirect(http.StatusFound, original)
}
