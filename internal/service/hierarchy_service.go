package service

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"

	"project/internal/models"
	"project/internal/repository"
)

type HierarchyService struct {
	repo       *repository.HierarchyRepository
	parser     *ExcelParser
	operations *HierarchyOperations
}

func NewHierarchyService(repo *repository.HierarchyRepository) *HierarchyService {
	service := &HierarchyService{repo: repo}
	service.parser = NewExcelParser(service)
	service.operations = NewHierarchyOperations(service)
	return service
}

func (s *HierarchyService) UploadExcel(file multipart.File) error {
	rows, err := s.parser.ParseExcelFile(file)
	if err != nil {
		return err
	}

	for _, row := range rows {
		if err := s.operations.ProcessUploadRow(row); err != nil {
			return err
		}
	}

	return nil
}

func (s *HierarchyService) GetAll() ([]models.HierarchyNode, error) {
	return s.repo.GetAll()
}

func (s *HierarchyService) CreateSingleNode(name, nodeType string, parentID *uuid.UUID, extraData map[string]string) (*models.HierarchyNode, error) {
	return s.operations.CreateNode(name, nodeType, parentID, extraData)
}

func (s *HierarchyService) GetHierarchyByNodeID(nodeID uuid.UUID) (*models.HierarchyNode, error) {
	return s.repo.GetByID(nodeID)
}

func (s *HierarchyService) GetSubTree(nodeID uuid.UUID) ([]models.TreeNode, error) {
	return s.operations.GetNodeSubTree(nodeID)
}

func (s *HierarchyService) GetFullHierarchy() ([]models.TreeNode, error) {
	return s.operations.GetFullHierarchyTree()
}

func (s *HierarchyService) findOrCreate(name, nodeType string, parentID *uuid.UUID, extra datatypes.JSON) (*models.HierarchyNode, error) {
	name = strings.TrimSpace(name)
	nodeType = strings.TrimSpace(nodeType)

	existing, err := s.repo.FindByNameAndParent(name, nodeType, parentID)
	if err == nil && existing != nil {
		return existing, nil
	}

	if err != nil && !repository.IsNotFound(err) {
		return nil, err
	}

	node := &models.HierarchyNode{
		ID:        uuid.New(),
		Name:      name,
		NodeType:  nodeType,
		ParentID:  parentID,
		ExtraData: extra,
	}

	return node, s.repo.Create(node)
}

func (s *HierarchyService) ExportToExcel() (*excelize.File, error) {
	nodes, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetName := "Hierarchy"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	// Set headers
	headers := []string{"ID", "Name", "Node Type", "Parent ID", "Extra Data", "Created At", "Updated At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Add data
	for i, node := range nodes {
		row := i + 2 // Start from row 2
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), node.ID.String())
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), node.Name)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), node.NodeType)
		if node.ParentID != nil {
			f.SetCellValue(sheetName, "D"+strconv.Itoa(row), node.ParentID.String())
		}
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), string(node.ExtraData))
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), node.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), node.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	return f, nil
}
