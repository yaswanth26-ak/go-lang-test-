const BASE_URL = "http://localhost:8080/api/v1/hierarchy";

export async function uploadHierarchy(file) {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(`${BASE_URL}/upload`, {
    method: "POST",
    body: formData,
  });

  return response.json();
}

export async function getHierarchy() {
  const response = await fetch(`${BASE_URL}/full`);
  return response.json();
}

export async function getNodeById(nodeId) {
  const response = await fetch(`${BASE_URL}/${nodeId}`);
  return response.json();
}

export async function getSubTree(nodeId) {
  const response = await fetch(`${BASE_URL}/${nodeId}/subtree`);
  return response.json();
}

export async function createNode(nodeData) {
  const response = await fetch(BASE_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(nodeData),
  });

  return response.json();
}

export async function downloadExcel() {
  const response = await fetch(`${BASE_URL}/download`);
  if (!response.ok) {
    throw new Error("Failed to download Excel file");
  }
  return response.blob();
}