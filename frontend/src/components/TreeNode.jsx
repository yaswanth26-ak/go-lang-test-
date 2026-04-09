export default function TreeNode({ node }) {
  if (!node) return null;

  return (
    <li className="tree-node">
      <div className="tree-node-name">
        {node.name}
      </div>
      
      <span className="tree-node-type">{node.node_type}</span>

      {node.extra_data && Object.keys(node.extra_data).length > 0 && (
        <div className="tree-node-meta">
          {Object.entries(node.extra_data).map(([key, value]) => (
            <div key={key}>
              {key}: {value}
            </div>
          ))}
        </div>
      )}

      {node.children && node.children.length > 0 && (
        <ul className="tree-node-nested">
          {node.children.map((child) => (
            <TreeNode key={child.id} node={child} />
          ))}
        </ul>
      )}
    </li>
  );
}