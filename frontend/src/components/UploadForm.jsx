import { useState } from "react";
import { uploadHierarchy } from "../api/hierarchyApi";

export default function UploadForm({ onUploadSuccess, setAlert }) {
  const [file, setFile] = useState(null);

  const handleUpload = async () => {
    if (!file) {
      setAlert("Please select a file first.", "error");
      return;
    }

    try {
      const result = await uploadHierarchy(file);
      if (result.message) {
        setAlert(result.message, "success");
        setFile(null);
        onUploadSuccess();
      } else {
        setAlert(result.error || "Upload failed.", "error");
      }
    } catch (error) {
      setAlert("Backend connection failed.", "error");
    }
  };

  return (
    <div className="control-card">
      <h3>Upload Hierarchy File</h3>
      
      <div className="form-group">
        <label>Pick Excel or CSV File</label>
        <input
          type="file"
          accept=".xlsx,.xls,.csv"
          onChange={(e) => setFile(e.target.files[0])}
        />
      </div>

      {file && (
        <p style={{ fontSize: "12px", color: "#666", marginBottom: "10px" }}>
          Selected: <strong>{file.name}</strong>
        </p>
      )}

      <button
        onClick={handleUpload}
        className="primary-button"
        disabled={!file}
      >
        Upload
      </button>
    </div>
  );
}