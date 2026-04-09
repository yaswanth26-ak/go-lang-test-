package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"project/internal/models"
)

type HierarchyRepository struct {
	db *gorm.DB
}

func NewHierarchyRepository(db *gorm.DB) *HierarchyRepository {
	return &HierarchyRepository{db: db}
}

func (r *HierarchyRepository) Create(node *models.HierarchyNode) error {
	return r.db.Create(node).Error
}

func (r *HierarchyRepository) GetAll() ([]models.HierarchyNode, error) {
	var nodes []models.HierarchyNode
	err := r.db.Order("created_at asc").Find(&nodes).Error
	return nodes, err
}

func (r *HierarchyRepository) FindByNameAndParent(name, nodeType string, parentID *uuid.UUID) (*models.HierarchyNode, error) {
	var node models.HierarchyNode

	name = strings.TrimSpace(name)
	nodeType = strings.TrimSpace(nodeType)

	query := r.db.
		Where("LOWER(name) = LOWER(?)", name).
		Where("LOWER(node_type) = LOWER(?)", nodeType)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.First(&node).Error
	return &node, err
}

func (r *HierarchyRepository) GetByID(id uuid.UUID) (*models.HierarchyNode, error) {
	var node models.HierarchyNode
	err := r.db.First(&node, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *HierarchyRepository) GetDB() *gorm.DB {
	return r.db
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
