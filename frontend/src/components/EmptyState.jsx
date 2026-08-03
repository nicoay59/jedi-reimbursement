import AppIcon from "./AppIcon";

export default function EmptyState({
  title,
  description,
  action,
  icon = "file"
}) {
  return (
    <div className="empty-state text-center px-3">
      <div className="empty-state-icon mx-auto mb-3">
        <AppIcon name={icon} size={28} />
      </div>
      <h2 className="h5 fw-bold">{title}</h2>
      <p className="text-secondary mb-4">{description}</p>
      {action}
    </div>
  );
}
