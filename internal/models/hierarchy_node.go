package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type HierarchyNode struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	NodeType  string         `gorm:"type:varchar(100);not null" json:"node_type"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	ExtraData datatypes.JSON `gorm:"type:jsonb" json:"extra_data,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (HierarchyNode) TableName() string {
	return "hierarchy_nodes"
}