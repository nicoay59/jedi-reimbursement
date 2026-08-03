import { shortDate } from "../utils/formatters";

export default function TrendChart({ items }) {
  if (!items || items.length === 0) {
    return (
      <div className="text-center text-secondary py-5">
        Belum ada data pada periode ini.
      </div>
    );
  }

  const maxValue = Math.max(
    1,
    ...items.map((item) => item.total_claims)
  );

  return (
    <div className="trend-chart" role="img" aria-label="Tren klaim harian">
      {items.map((item) => {
        const height = Math.max(
          8,
          Math.round((item.total_claims / maxValue) * 100)
        );

        return (
          <div className="trend-chart-column" key={item.date}>
            <div className="trend-chart-value">
              {item.total_claims}
            </div>

            <div className="trend-chart-bar-track">
              <div
                className="trend-chart-bar"
                style={{ height: `${height}%` }}
                title={`${item.total_claims} klaim`}
              />
            </div>

            <div className="trend-chart-label">
              {shortDate(item.date)}
            </div>
          </div>
        );
      })}
    </div>
  );
}
