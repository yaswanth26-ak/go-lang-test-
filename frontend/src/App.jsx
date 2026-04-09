import { useState } from "react";
import UploadForm from "./components/UploadForm";
import TreeNode from "./components/TreeNode";
import NodeForm from "./components/NodeForm";
import { getHierarchy, getSubTree, downloadExcel } from "./api/hierarchyApi";
import "./App.css";

export default function App() {
  const [data, setData] = useState([]);
  const [selectedNodeId, setSelectedNodeId] = useState("");
  const [message, setMessage] = useState({ type: "", text: "" });
  const [loading, setLoading] = useState(false);

  const setAlert = (text, type = "success") => {
    setMessage({ text, type });
    setTimeout(() => setMessage({ type: "", text: "" }), 4000);
  };

  const loadFull = async () => {
    setLoading(true);
    setAlert("Loading...", "info");
    try {
      const result = await getHierarchy();
      setData(result);
      setAlert("Loaded successfully!", "success");
    } catch (error) {
      setAlert("Failed to load.", "error");
    } finally {
      setLoading(false);
    }
  };

  const loadSub = async () => {
    if (!selectedNodeId.trim()) {
      setAlert("Please enter node ID.", "error");
      return;
    }

    setLoading(true);
    setAlert("Loading subtree...", "info");
    try {
      const result = await getSubTree(selectedNodeId);
      setData(result);
      setAlert("Subtree loaded!", "success");
    } catch (error) {
      setAlert("Failed to load subtree.", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleNodeCreated = () => {
    loadFull();
  };

  const downloadHierarchy = async () => {
    setAlert("Downloading...", "info");
    try {
      const blob = await downloadExcel();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "hierarchy.xlsx";
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
      setAlert("Downloaded successfully!", "success");
    } catch (error) {
      setAlert("Failed to download.", "error");
    }
  };

  return (
    <div className="app-shell">
      <header className="app-header">
        <h1>Hierarchy Manager</h1>
        <p className="subtitle">Upload files and manage your asset hierarchy</p>
      </header>

      <section className="controls-grid">
        <UploadForm onUploadSuccess={loadFull} setAlert={setAlert} />
        <NodeForm onNodeCreated={handleNodeCreated} setAlert={setAlert} />
      </section>

      <section className="actions-card">
        <div className="action-row">
          <button className="primary-button" onClick={loadFull} disabled={loading}>
            {loading ? "Loading..." : "Load All"}
          </button>
          <button className="secondary-button" onClick={downloadHierarchy} disabled={loading}>
            Download Excel
          </button>
          <div className="subtree-input">
            <input
              type="text"
              placeholder="Enter node ID..."
              value={selectedNodeId}
              onChange={(e) => setSelectedNodeId(e.target.value)}
              disabled={loading}
            />
            <button className="secondary-button" onClick={loadSub} disabled={loading}>
              Search
            </button>
          </div>
        </div>
        <div className="meta-row">
          <span>{data.length} nodes</span>
          <span>{loading ? "Working..." : "Ready"}</span>
        </div>
      </section>

      {message.text && (
        <div className={`status-banner ${message.type}`}>
          {message.text}
        </div>
      )}

      <main className="tree-panel">
        <div className="panel-header">
          <h2>Tree View</h2>
          <p>Your hierarchy structure</p>
        </div>

        {data.length > 0 ? (
          <ul className="tree-list">
            {data.map((node) => (
              <TreeNode key={node.id} node={node} />
            ))}
          </ul>
        ) : (
          <div className="empty-state">
            <p>No data yet</p>
            <button className="primary-button" onClick={loadFull} disabled={loading}>
              Load Hierarchy
            </button>
          </div>
        )}
      </main>
    </div>
  );
}
