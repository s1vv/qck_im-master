package handlers

import (
	"database/sql"
	"net/http"
	"qckim-backend/internal/logger"
	"qckim-backend/internal/models"
	"qckim-backend/internal/repository"
	"qckim-backend/internal/services"
	"qckim-backend/utils/valid"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QckLinkHandler struct {
	QckLinkService *services.QckLinkRepo
}

func NewQckLinkHandler(db *repository.QckRepo) *QckLinkHandler {
	return &QckLinkHandler{QckLinkService: services.NewQckLink(db.GetDB())}
}

func (h *QckLinkHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("api/qck/qck-link", h.GetQckLinkData)
	r.GET("api/qck/qck-links", h.GetQckAllLinks)
	r.GET("api/qck/shared-data-link", h.GetSharedDataLink)
	r.POST("api/qck/update-data-link", h.UpdateDataLink)
	r.POST("api/qck/activate-link", h.ActivateLink)
	r.POST("api/qck/remove-link-description", h.RemoveDescription)
}

func (h *QckLinkHandler) RemoveDescription(c *gin.Context) {
	var req models.UpdateDataLink

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.QckLinkService.UpdateDataLink("", "", req.QckLink)
	if err != nil {
		logger.Error("(h *QckLinkHandler) RemoveDescription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Remove description"})
}

func (h *QckLinkHandler) ActivateLink(c *gin.Context) {
	var req models.ActivationData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		logger.Error("ActivateLink", "!ok", ok)
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	userID64, err := strconv.ParseInt(userID.(string), 10, 64)
	if err != nil {
		logger.Error("ActivateLink", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	err = h.QckLinkService.ActivateLink(userID64, req.QckLink, req.Password)
	if err != nil {
		logger.Error("h.QckLinkService.ActivateLink", "error", err)
	}

	if !ok {
		logger.Error("ActivateLink", "!ok", ok)
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link activated"})
}

func (h *QckLinkHandler) UpdateDataLink(c *gin.Context) {
	var req models.UpdateDataLink

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	err := h.QckLinkService.UpdateDataLink(req.Name, req.Description, req.QckLink)
	if err != nil {
		logger.Error("h.QckLinkService.UpdateDataLink", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})

}
func (h *QckLinkHandler) GetSharedDataLink(c *gin.Context) {
	qckLink := c.DefaultQuery("link", "")
	if !valid.ValidateQckLink(qckLink) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid qck link"})
		return
	}
	_, d, err := h.QckLinkService.GetQckLinkDescription(qckLink)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "qck-link not found"})
		return
	} else if err != nil {
		logger.Error("GerQckLinkData", "message", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	if len(d) == 0 || d == "qck.im" {
		c.JSON(http.StatusNotFound, gin.H{"message": "qck-link not found"})
	}
	c.JSON(http.StatusOK, gin.H{"description": d})

}

func (h *QckLinkHandler) GetQckLinkData(c *gin.Context) {
	qckLink := c.DefaultQuery("link", "")
	if !valid.ValidateQckLink(qckLink) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid qck link"})
		return
	}
	name, d, err := h.QckLinkService.GetQckLinkDescription(qckLink)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "qck-link not found"})
		return
	} else if err != nil {
		logger.Error("GerQckLinkData", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"description": d, "name": name})
}

func (h *QckLinkHandler) GetQckAllLinks(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		logger.Error("GetQckAllLinks", "!ok", ok)
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	userID64, err := strconv.ParseInt(userID.(string), 10, 64)
	if err != nil {
		logger.Error("GetQckAllLinks", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	qckLinksData, err := h.QckLinkService.GetAllUserLinks(userID64)
	if err != nil {
		logger.Error("GetQckAllLinks", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": qckLinksData})
}
