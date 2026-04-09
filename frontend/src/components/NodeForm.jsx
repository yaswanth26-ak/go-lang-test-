import { useState } from "react";
import { createNode } from "../api/hierarchyApi";

export default function NodeForm({ onNodeCreated, setAlert }) {
  const [name, setName] = useState("");
  const [nodeType, setNodeType] = useState("");
  const [parentId, setParentId] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!name.trim() || !nodeType.trim()) {
      setAlert("Name and Node Type are required.", "error");
      return;
    }

    try {
      const nodeData = {
        name: name,
        node_type: nodeType,
        parent_id: parentId.trim() || null,
        extra_data: {},
      };

      const result = await createNode(nodeData);

      if (result.id) {
        setAlert("Node created!", "success");
        setName("");
        setNodeType("");
        setParentId("");
        onNodeCreated();
      } else {
        setAlert(result.error || "Failed to create node.", "error");
      }
    } catch (error) {
      setAlert("Connection failed.", "error");
    }
  };

  return (
    <div className="control-card">
      <h3>Create New Node</h3>
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Node Name</label>
          <input
            type="text"
            placeholder="e.g., Department A"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>

        <div className="form-group">
          <label>Node Type</label>
          <input
            type="text"
            placeholder="e.g., Department"
            value={nodeType}
            onChange={(e) => setNodeType(e.target.value)}
            required
          />
        </div>

        <div className="form-group">
          <label>Parent ID (optional)</label>
          <input
            type="text"
            placeholder="Leave empty for root"
            value={parentId}
            onChange={(e) => setParentId(e.target.value)}
          />
        </div>

        <button type="submit" className="primary-button">
          Create
        </button>
      </form>
    </div>
  );
}