package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"project/internal/models"
	"project/internal/service"
)

type HierarchyHandler struct {
	service *service.HierarchyService
}

func NewHierarchyHandler(service *service.HierarchyService) *HierarchyHandler {
	return &HierarchyHandler{service: service}
}

// UploadExcel godoc
// @Summary Upload Excel file
// @Description Upload hierarchy Excel file and store in DB
// @Tags hierarchy
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel File"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/hierarchy/upload [post]
func (h *HierarchyHandler) UploadExcel(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			h.respondWithError(c, http.StatusInternalServerError, "internal server error during file processing")
		}
	}()
	
	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, "file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, "unable to open file")
		return
	}
	defer file.Close()

	if err := h.service.UploadExcel(file); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "excel uploaded successfully"})
}

// GetAll godoc
// @Summary Get all hierarchy data
// @Description Fetch all hierarchy nodes
// @Tags hierarchy
// @Produce json
// @Success 200 {array} models.TreeNode
// @Failure 500 {object} map[string]string
// @Router /api/v1/hierarchy [get]
func (h *HierarchyHandler) GetAll(c *gin.Context) {
	nodes, err := h.service.GetAll()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "failed to fetch data")
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// DownloadExcel godoc
// @Summary Download hierarchy data as Excel
// @Description Export all hierarchy data to Excel file
// @Tags hierarchy
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} binary
// @Failure 500 {object} map[string]string
// @Router /api/v1/hierarchy/download [get]
func (h *HierarchyHandler) DownloadExcel(c *gin.Context) {
	file, err := h.service.ExportToExcel()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "failed to generate excel file")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=hierarchy.xlsx")
	c.Header("Cache-Control", "no-cache")

	file.Write(c.Writer)
}

// CreateSingleNode godoc
// @Summary Create a single hierarchy node
// @Description Create a new hierarchy node manually
// @Tags hierarchy
// @Accept json
// @Produce json
// @Param node body map[string]interface{} true "Node data"
// @Success 201 {object} models.TreeNode
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/hierarchy [post]
func (h *HierarchyHandler) CreateSingleNode(c *gin.Context) {
	var req struct {
		Name      string            `json:"name" binding:"required"`
		NodeType  string            `json:"node_type" binding:"required"`
		ParentID  *string           `json:"parent_id,omitempty"`
		ExtraData map[string]string `json:"extra_data,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	parentID, err := h.parseParentID(req.ParentID)
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	node, err := h.service.CreateSingleNode(req.Name, req.NodeType, parentID, req.ExtraData)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, node)
}

// GetHierarchyByNodeID godoc
// @Summary Get hierarchy node by ID
// @Description Fetch a specific hierarchy node and its details
// @Tags hierarchy
// @Produce json
// @Param id path string true "Node ID"
// @Success 200 {object} models.TreeNode
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/hierarchy/{id} [get]
func (h *HierarchyHandler) GetHierarchyByNodeID(c *gin.Context) {
	nodeID, err := h.parseNodeID(c.Param("id"))
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	node, err := h.service.GetHierarchyByNodeID(nodeID)
	if err != nil {
		h.respondWithError(c, http.StatusNotFound, "node not found")
		return
	}

	c.JSON(http.StatusOK, node)
}

// GetSubTree godoc
// @Summary Get subtree for a node
// @Description Fetch partial hierarchy (subtree) starting from a specific node
// @Tags hierarchy
// @Produce json
// @Param id path string true "Node ID"
// @Success 200 {array} models.TreeNode
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/hierarchy/{id}/subtree [get]
func (h *HierarchyHandler) GetSubTree(c *gin.Context) {
	nodeID, err := h.parseNodeID(c.Param("id"))
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	nodes, err := h.service.GetSubTree(nodeID)
	if err != nil {
		h.respondWithError(c, http.StatusNotFound, "node not found or failed to fetch subtree")
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// GetFullHierarchy godoc
// @Summary Get full hierarchy
// @Description Fetch complete hierarchy in tree structure
// @Tags hierarchy
// @Produce json
// @Success 200 {array} models.TreeNode
// @Failure 500 {object} map[string]string
// @Router /api/v1/hierarchy/full [get]
func (h *HierarchyHandler) GetFullHierarchy(c *gin.Context) {
	nodes, err := h.service.GetFullHierarchy()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "failed to fetch full hierarchy")
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// Helper methods
func (h *HierarchyHandler) parseNodeID(idStr string) (uuid.UUID, error) {
	return uuid.Parse(idStr)
}

func (h *HierarchyHandler) parseParentID(parentIDStr *string) (*uuid.UUID, error) {
	if parentIDStr == nil {
		return nil, nil
	}

	parsedID, err := uuid.Parse(*parentIDStr)
	if err != nil {
		return nil, err
	}
	return &parsedID, nil
}

func (h *HierarchyHandler) respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func (h *HierarchyHandler) handleServiceError(c *gin.Context, err error) {
	if validationErr, ok := err.(service.ValidationError); ok {
		h.respondWithError(c, http.StatusBadRequest, validationErr.Error())
		return
	}
	h.respondWithError(c, http.StatusInternalServerError, "internal server error")
}

// Swagger helper to ensure tree model is included
var _ = models.TreeNode{}
