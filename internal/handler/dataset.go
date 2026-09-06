package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lattiq-bhuvan/access-desk/internal/store"
)

type DatasetHandler struct {
	store *store.Store
}

func NewDatasetHandler(s *store.Store) *DatasetHandler {
	return &DatasetHandler{store: s}
}

func (h *DatasetHandler) List(c *gin.Context) {
	datasets, err := h.store.ListDatasets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list datasets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"datasets": datasets})
}