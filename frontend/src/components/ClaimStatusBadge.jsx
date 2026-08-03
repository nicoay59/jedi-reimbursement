const statusConfig = {
  PENDING: { label: "Menunggu", className: "status-pending" },
  APPROVED: { label: "Disetujui", className: "status-approved" },
  REJECTED: { label: "Ditolak", className: "status-rejected" }
};

export default function ClaimStatusBadge({ status }) {
  const config = statusConfig[status] || {
    label: status || "-",
    className: "status-default"
  };

  return (
    <span className={`status-badge ${config.className}`}>
      <span className="status-dot" aria-hidden="true" />
      {config.label}
    </span>
  );
}
