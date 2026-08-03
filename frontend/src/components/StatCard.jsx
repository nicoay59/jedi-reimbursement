import AppIcon from "./AppIcon";

export default function StatCard({
  label,
  value,
  helper,
  tone = "primary",
  icon = "total"
}) {
  return (
    <article className={`card app-card stat-card stat-card-${tone} h-100`}>
      <div className="card-body p-4">
        <div className="d-flex justify-content-between align-items-start gap-3">
          <div>
            <p className="stat-card-label mb-2">{label}</p>
            <div className="stat-card-value">{value}</div>
          </div>
          <span className="stat-card-icon">
            <AppIcon name={icon} size={22} />
          </span>
        </div>
        {helper && <p className="stat-card-helper mt-3 mb-0">{helper}</p>}
      </div>
    </article>
  );
}
