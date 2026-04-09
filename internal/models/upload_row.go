package models

type UploadRow struct {
	Levels    []LevelData
	ExtraData map[string]string
}

type LevelData struct {
	Header string
	Value  string
}

// TreeNode represents a node in the hierarchy tree
type TreeNode struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	NodeType  string                 `json:"node_type"`
	ParentID  *string                `json:"parent_id,omitempty"`
	ExtraData map[string]interface{} `json:"extra_data,omitempty"`
	Children  []TreeNode             `json:"children"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}