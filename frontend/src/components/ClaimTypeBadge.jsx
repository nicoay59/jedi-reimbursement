import AppIcon from "./AppIcon";

const typeConfig = {
  PARKING: { label: "Parkir", className: "type-parking", icon: "parking" },
  OVERTIME: { label: "Lembur", className: "type-overtime", icon: "overtime" }
};

export default function ClaimTypeBadge({ type }) {
  const config = typeConfig[type] || {
    label: type || "-",
    className: "type-default",
    icon: "file"
  };

  return (
    <span className={`claim-type-badge ${config.className}`}>
      <AppIcon name={config.icon} size={14} />
      {config.label}
    </span>
  );
}
