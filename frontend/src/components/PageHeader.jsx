export default function PageHeader({
  eyebrow,
  title,
  description,
  actions
}) {
  return (
    <header className="page-header mb-4">
      <div className="d-flex flex-column flex-lg-row justify-content-between align-items-lg-end gap-3">
        <div className="page-header-copy">
          {eyebrow && <p className="page-eyebrow mb-2">{eyebrow}</p>}
          <h1 className="page-title mb-2">{title}</h1>
          {description && (
            <p className="page-description mb-0">{description}</p>
          )}
        </div>

        {actions && <div className="page-actions d-flex flex-wrap gap-2">{actions}</div>}
      </div>
    </header>
  );
}
