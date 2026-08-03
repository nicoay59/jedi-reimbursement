import { useEffect, useState } from "react";
import ClaimStatusBadge from "../components/ClaimStatusBadge";
import ClaimTypeBadge from "../components/ClaimTypeBadge";
import DashboardLayout from "../components/DashboardLayout";
import PageHeader from "../components/PageHeader";
import StatCard from "../components/StatCard";
import AppIcon from "../components/AppIcon";
import TrendChart from "../components/TrendChart";
import { getAdminDashboard } from "../services/reportService";
import {
  currentMonthRange,
  formatDate,
  formatDateRange,
  formatDuration,
  formatRupiah
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

const initialSummary = {
  total_claims: 0,
  pending_claims: 0,
  approved_claims: 0,
  rejected_claims: 0,
  parking_claims: 0,
  overtime_claims: 0,
  total_parking_amount: 0,
  approved_parking_amount: 0,
  total_overtime_hours: 0,
  approved_overtime_hours: 0
};

export default function AdminDashboardPage() {
  const defaultRange = currentMonthRange();
  const [draftPeriod, setDraftPeriod] = useState(defaultRange);
  const [period, setPeriod] = useState(defaultRange);
  const [state, setState] = useState({
    loading: true,
    error: "",
    summary: initialSummary,
    trend: [],
    recent: []
  });

  useEffect(() => {
    let active = true;

    async function loadDashboard() {
      setState((current) => ({
        ...current,
        loading: true,
        error: ""
      }));

      try {
        const response = await getAdminDashboard(period);

        if (!active) return;

        setState({
          loading: false,
          error: "",
          summary: response.data.summary,
          trend: response.data.trend,
          recent: response.data.recent
        });
      } catch (error) {
        if (!active) return;

        setState((current) => ({
          ...current,
          loading: false,
          error: error.message
        }));
      }
    }

    loadDashboard();

    return () => {
      active = false;
    };
  }, [period]);

  function updatePeriod(event) {
    const { name, value } = event.target;
    setDraftPeriod((current) => ({
      ...current,
      [name]: value
    }));
  }

  function applyPeriod(event) {
    event.preventDefault();

    if (
      !draftPeriod.startDate ||
      !draftPeriod.endDate ||
      draftPeriod.endDate < draftPeriod.startDate
    ) {
      setState((current) => ({
        ...current,
        error:
          "Tanggal selesai harus sama dengan atau setelah tanggal mulai."
      }));
      return;
    }

    setPeriod({ ...draftPeriod });
  }

  const summary = state.summary;

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Dashboard Administrator"
        title="Ringkasan reimbursement"
        description="Pantau jumlah, nilai, durasi, dan status klaim berdasarkan periode."
        actions={
          <>
            <button
              type="button"
              className="btn btn-outline-primary"
              onClick={() => navigate("/admin/claims")}
            >
              <AppIcon name="review" size={17} className="me-2" />
              Pemeriksaan Klaim
            </button>
            <button
              type="button"
              className="btn btn-primary"
              onClick={() => navigate("/admin/reports")}
            >
              <AppIcon name="report" size={17} className="me-2" />
              Buka Laporan
            </button>
          </>
        }
      />

      <form
        className="card app-card mb-4"
        onSubmit={applyPeriod}
      >
        <div className="card-body p-4">
          <div className="row g-3 align-items-end">
            <div className="col-md-5">
              <label htmlFor="startDate" className="form-label">
                Tanggal mulai
              </label>
              <input
                id="startDate"
                name="startDate"
                type="date"
                className="form-control"
                value={draftPeriod.startDate}
                max={draftPeriod.endDate}
                onChange={updatePeriod}
              />
            </div>

            <div className="col-md-5">
              <label htmlFor="endDate" className="form-label">
                Tanggal selesai
              </label>
              <input
                id="endDate"
                name="endDate"
                type="date"
                className="form-control"
                value={draftPeriod.endDate}
                min={draftPeriod.startDate}
                onChange={updatePeriod}
              />
            </div>

            <div className="col-md-2">
              <button
                type="submit"
                className="btn btn-primary w-100"
                disabled={state.loading}
              >
                {state.loading ? "Memuat..." : "Terapkan"}
              </button>
            </div>
          </div>
        </div>
      </form>

      {state.error && (
        <div className="alert alert-danger">{state.error}</div>
      )}

      <div className="row g-4 mb-4">
        <div className="col-sm-6 col-xl-3">
          <StatCard
            label="Total Klaim"
            value={summary.total_claims}
            helper={`${summary.parking_claims} parkir · ${summary.overtime_claims} lembur`}
            icon="total"
          />
        </div>

        <div className="col-sm-6 col-xl-3">
          <StatCard
            label="Menunggu"
            value={summary.pending_claims}
            helper="Perlu diperiksa administrator"
            tone="warning"
            icon="pending"
          />
        </div>

        <div className="col-sm-6 col-xl-3">
          <StatCard
            label="Disetujui"
            value={summary.approved_claims}
            helper="Pengajuan yang telah diterima"
            tone="success"
            icon="check"
          />
        </div>

        <div className="col-sm-6 col-xl-3">
          <StatCard
            label="Ditolak"
            value={summary.rejected_claims}
            helper="Pengajuan yang tidak disetujui"
            tone="danger"
            icon="close"
          />
        </div>
      </div>

      <div className="row g-4 mb-4">
        <div className="col-lg-6">
          <div className="card app-card h-100">
            <div className="card-body p-4">
              <p className="text-primary fw-semibold mb-1">
                Klaim Parkir
              </p>
              <h2 className="h5 mb-4">Ringkasan nominal</h2>

              <div className="report-metric-grid">
                <div className="report-metric">
                  <span className="text-secondary small">
                    Total diajukan
                  </span>
                  <strong>
                    {formatRupiah(summary.total_parking_amount)}
                  </strong>
                </div>

                <div className="report-metric">
                  <span className="text-secondary small">
                    Total disetujui
                  </span>
                  <strong>
                    {formatRupiah(
                      summary.approved_parking_amount
                    )}
                  </strong>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="col-lg-6">
          <div className="card app-card h-100">
            <div className="card-body p-4">
              <p className="text-primary fw-semibold mb-1">
                Klaim Lembur
              </p>
              <h2 className="h5 mb-4">Ringkasan durasi</h2>

              <div className="report-metric-grid">
                <div className="report-metric">
                  <span className="text-secondary small">
                    Total diajukan
                  </span>
                  <strong>
                    {formatDuration(
                      summary.total_overtime_hours
                    )}
                  </strong>
                </div>

                <div className="report-metric">
                  <span className="text-secondary small">
                    Total disetujui
                  </span>
                  <strong>
                    {formatDuration(
                      summary.approved_overtime_hours
                    )}
                  </strong>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="row g-4">
        <div className="col-xl-7">
          <div className="card app-card h-100">
            <div className="card-body p-4">
              <p className="text-primary fw-semibold mb-1">
                Tren Harian
              </p>
              <h2 className="h5 mb-4">
                Jumlah pengajuan per tanggal
              </h2>

              <TrendChart items={state.trend} />
            </div>
          </div>
        </div>

        <div className="col-xl-5">
          <div className="card app-card h-100">
            <div className="card-body p-4">
              <div className="d-flex justify-content-between align-items-start gap-3 mb-4">
                <div>
                  <p className="text-primary fw-semibold mb-1">
                    Klaim Terbaru
                  </p>
                  <h2 className="h5 mb-0">
                    Pengajuan pada periode
                  </h2>
                </div>

                <button
                  type="button"
                  className="btn btn-outline-primary btn-sm"
                  onClick={() => navigate("/admin/reports")}
                >
                  Semua
                </button>
              </div>

              {state.recent.length === 0 ? (
                <p className="text-secondary mb-0">
                  Belum ada pengajuan pada periode ini.
                </p>
              ) : (
                <div className="recent-claim-list">
                  {state.recent.map((claim) => (
                    <button
                      type="button"
                      className="recent-claim-item"
                      key={`${claim.claim_type}-${claim.claim_id}`}
                      onClick={() =>
                        navigate(
                          `/admin/claims/${claim.claim_type}/${claim.claim_id}`
                        )
                      }
                    >
                      <div className="d-flex justify-content-between gap-3">
                        <div className="min-w-0">
                          <div className="d-flex flex-wrap gap-2 mb-2">
                            <ClaimTypeBadge
                              type={claim.claim_type}
                            />
                            <ClaimStatusBadge
                              status={claim.status}
                            />
                          </div>
                          <strong className="d-block text-truncate">
                            {claim.employee_name}
                          </strong>
                          <span className="small text-secondary">
                            {claim.claim_type === "PARKING"
                              ? formatDateRange(claim.claim_date, claim.claim_end_date)
                              : formatDate(claim.claim_date)}
                          </span>
                        </div>

                        <div className="text-end">
                          <span className="small text-secondary">
                            {claim.claim_type === "PARKING"
                              ? formatRupiah(claim.amount)
                              : formatDuration(
                                  claim.duration_hours
                                )}
                          </span>
                        </div>
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
