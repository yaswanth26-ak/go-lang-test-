package service

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"project/internal/models"
)

type HierarchyOperations struct {
	service *HierarchyService
}

func NewHierarchyOperations(service *HierarchyService) *HierarchyOperations {
	return &HierarchyOperations{service: service}
}

func (ho *HierarchyOperations) ProcessUploadRow(row models.UploadRow) error {
	var parentID *uuid.UUID

	for index, level := range row.Levels {
		extra := buildExtraData(row.ExtraData)
		if index > 0 {
			extra = datatypes.JSON([]byte("{}"))
		}

		node, err := ho.service.findOrCreate(level.Value, level.Header, parentID, extra)
		if err != nil {
			return err
		}

		parentID = &node.ID
	}

	return nil
}

func (ho *HierarchyOperations) CreateNode(name, nodeType string, parentID *uuid.UUID, extraData map[string]string) (*models.HierarchyNode, error) {
	if err := ho.validateNodeInput(name, nodeType); err != nil {
		return nil, err
	}

	extra := buildExtraData(extraData)
	return ho.service.findOrCreate(name, nodeType, parentID, extra)
}

func (ho *HierarchyOperations) validateNodeInput(name, nodeType string) error {
	if strings.TrimSpace(name) == "" {
		return NewValidationError("name is required")
	}
	if strings.TrimSpace(nodeType) == "" {
		return NewValidationError("node_type is required")
	}
	return nil
}

func (ho *HierarchyOperations) GetNodeSubTree(nodeID uuid.UUID) ([]models.TreeNode, error) {
	root, err := ho.service.repo.GetByID(nodeID)
	if err != nil {
		return nil, err
	}

	return ho.buildTreeStructure([]models.HierarchyNode{*root})
}

func (ho *HierarchyOperations) GetFullHierarchyTree() ([]models.TreeNode, error) {
	allNodes, err := ho.getAllNodes()
	if err != nil {
		return nil, err
	}

	return ho.buildTreeStructure(allNodes)
}

func (ho *HierarchyOperations) getRootNodes() ([]models.HierarchyNode, error) {
	var roots []models.HierarchyNode
	err := ho.service.repo.GetDB().Where("parent_id IS NULL").Find(&roots).Error
	return roots, err
}

func (ho *HierarchyOperations) getAllNodes() ([]models.HierarchyNode, error) {
	return ho.service.repo.GetAll()
}

func (ho *HierarchyOperations) buildTreeStructure(nodes []models.HierarchyNode) ([]models.TreeNode, error) {
	nodeMap := make(map[string]*models.TreeNode)

	// First pass: create all TreeNodes
	for _, node := range nodes {
		treeNode := ho.convertToTreeNode(node)
		nodeMap[node.ID.String()] = &treeNode
	}

	// Second pass: build parent-child relationships
	for _, node := range nodes {
		if node.ParentID != nil {
			if parent, exists := nodeMap[node.ParentID.String()]; exists {
				if child, exists := nodeMap[node.ID.String()]; exists {
					parent.Children = append(parent.Children, *child)
				}
			}
		}
	}

	// Collect root nodes
	var rootNodes []models.TreeNode
	for _, node := range nodes {
		if node.ParentID == nil {
			if root, exists := nodeMap[node.ID.String()]; exists {
				rootNodes = append(rootNodes, *root)
			}
		}
	}

	return rootNodes, nil
}

func (ho *HierarchyOperations) convertToTreeNode(node models.HierarchyNode) models.TreeNode {
	var extraData map[string]interface{}
	if len(node.ExtraData) > 0 {
		json.Unmarshal(node.ExtraData, &extraData)
	}

	var parentID *string
	if node.ParentID != nil {
		idStr := node.ParentID.String()
		parentID = &idStr
	}

	return models.TreeNode{
		ID:        node.ID.String(),
		Name:      node.Name,
		NodeType:  node.NodeType,
		ParentID:  parentID,
		ExtraData: extraData,
		Children:  nil, // Don't initialize, let append create the slice
		CreatedAt: node.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: node.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
