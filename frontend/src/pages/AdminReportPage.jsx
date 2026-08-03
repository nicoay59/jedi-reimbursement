import { useEffect, useState } from "react";
import ClaimStatusBadge from "../components/ClaimStatusBadge";
import ClaimTypeBadge from "../components/ClaimTypeBadge";
import DashboardLayout from "../components/DashboardLayout";
import EmptyState from "../components/EmptyState";
import LoadingScreen from "../components/LoadingScreen";
import PageHeader from "../components/PageHeader";
import {
  exportAdminReports,
  getAdminReports
} from "../services/reportService";
import {
  currentMonthRange,
  formatDate,
  formatDateRange,
  formatDuration,
  formatRupiah
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

const typeOptions = [
  { value: "ALL", label: "Semua jenis" },
  { value: "PARKING", label: "Parkir" },
  { value: "OVERTIME", label: "Lembur" }
];

const statusOptions = [
  { value: "ALL", label: "Semua status" },
  { value: "PENDING", label: "Menunggu" },
  { value: "APPROVED", label: "Disetujui" },
  { value: "REJECTED", label: "Ditolak" }
];

export default function AdminReportPage() {
  const defaultRange = currentMonthRange();
  const defaultFilters = {
    ...defaultRange,
    type: "ALL",
    status: "ALL"
  };

  const [draftFilters, setDraftFilters] =
    useState(defaultFilters);
  const [filters, setFilters] = useState(defaultFilters);
  const [page, setPage] = useState(1);
  const [state, setState] = useState({
    loading: true,
    exporting: false,
    error: "",
    items: [],
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_pages: 0
    }
  });

  useEffect(() => {
    let active = true;

    async function loadReports() {
      setState((current) => ({
        ...current,
        loading: true,
        error: ""
      }));

      try {
        const response = await getAdminReports({
          ...filters,
          page,
          limit: 10
        });

        if (!active) return;

        setState((current) => ({
          ...current,
          loading: false,
          error: "",
          items: response.data.items,
          pagination: response.data.pagination
        }));
      } catch (error) {
        if (!active) return;

        setState((current) => ({
          ...current,
          loading: false,
          error: error.message
        }));
      }
    }

    loadReports();

    return () => {
      active = false;
    };
  }, [filters, page]);

  function updateFilter(event) {
    const { name, value } = event.target;
    setDraftFilters((current) => ({
      ...current,
      [name]: value
    }));
  }

  function applyFilters(event) {
    event.preventDefault();

    if (
      !draftFilters.startDate ||
      !draftFilters.endDate ||
      draftFilters.endDate < draftFilters.startDate
    ) {
      setState((current) => ({
        ...current,
        error:
          "Tanggal selesai harus sama dengan atau setelah tanggal mulai."
      }));
      return;
    }

    setPage(1);
    setFilters({ ...draftFilters });
  }

  async function exportCSV() {
    setState((current) => ({
      ...current,
      exporting: true,
      error: ""
    }));

    try {
      const response = await exportAdminReports(filters);
      const url = URL.createObjectURL(response.blob);
      const link = document.createElement("a");

      link.href = url;
      link.download =
        `jedi-reimbursement-${filters.startDate}-${filters.endDate}.csv`;
      document.body.appendChild(link);
      link.click();
      link.remove();

      setTimeout(() => URL.revokeObjectURL(url), 1000);

      setState((current) => ({
        ...current,
        exporting: false
      }));
    } catch (error) {
      setState((current) => ({
        ...current,
        exporting: false,
        error: error.message
      }));
    }
  }

  if (state.loading && state.items.length === 0) {
    return (
      <DashboardLayout>
        <LoadingScreen message="Memuat laporan..." />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Administrator"
        title="Laporan reimbursement"
        description="Rekapitulasi klaim parkir dan lembur berdasarkan periode."
        actions={
          <>
            <button
              type="button"
              className="btn btn-outline-secondary"
              onClick={() => navigate("/admin/dashboard")}
            >
              Dashboard
            </button>

            <button
              type="button"
              className="btn btn-success"
              onClick={exportCSV}
              disabled={state.exporting}
            >
              {state.exporting
                ? "Membuat CSV..."
                : "Ekspor CSV"}
            </button>
          </>
        }
      />

      <form
        className="card app-card mb-4"
        onSubmit={applyFilters}
      >
        <div className="card-body p-4">
          <div className="row g-3 align-items-end">
            <div className="col-md-6 col-xl-3">
              <label htmlFor="startDate" className="form-label">
                Tanggal mulai
              </label>
              <input
                id="startDate"
                name="startDate"
                type="date"
                className="form-control"
                value={draftFilters.startDate}
                max={draftFilters.endDate}
                onChange={updateFilter}
              />
            </div>

            <div className="col-md-6 col-xl-3">
              <label htmlFor="endDate" className="form-label">
                Tanggal selesai
              </label>
              <input
                id="endDate"
                name="endDate"
                type="date"
                className="form-control"
                value={draftFilters.endDate}
                min={draftFilters.startDate}
                onChange={updateFilter}
              />
            </div>

            <div className="col-md-6 col-xl-2">
              <label htmlFor="type" className="form-label">
                Jenis
              </label>
              <select
                id="type"
                name="type"
                className="form-select"
                value={draftFilters.type}
                onChange={updateFilter}
              >
                {typeOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            <div className="col-md-6 col-xl-2">
              <label htmlFor="status" className="form-label">
                Status
              </label>
              <select
                id="status"
                name="status"
                className="form-select"
                value={draftFilters.status}
                onChange={updateFilter}
              >
                {statusOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            <div className="col-xl-2">
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

      <div className="card app-card">
        {state.items.length === 0 ? (
          <EmptyState
            title="Data laporan tidak ditemukan"
            description="Ubah periode atau filter untuk menampilkan data lain."
          />
        ) : (
          <>
            <div className="table-responsive">
              <table className="table align-middle mb-0">
                <thead className="table-light">
                  <tr>
                    <th className="ps-4">Jenis</th>
                    <th>Karyawan</th>
                    <th>Tanggal</th>
                    <th>Nilai / Durasi</th>
                    <th>Status</th>
                    <th>Pemeriksa</th>
                    <th className="text-end pe-4">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {state.items.map((claim) => (
                    <tr
                      key={`${claim.claim_type}-${claim.claim_id}`}
                    >
                      <td className="ps-4">
                        <ClaimTypeBadge type={claim.claim_type} />
                      </td>
                      <td>
                        <div className="fw-semibold">
                          {claim.employee_name}
                        </div>
                        <div className="small text-secondary">
                          {claim.employee_number}
                        </div>
                      </td>
                      <td>
                        {claim.claim_type === "PARKING"
                          ? formatDateRange(claim.claim_date, claim.claim_end_date)
                          : formatDate(claim.claim_date)}
                      </td>
                      <td>
                        {claim.claim_type === "PARKING"
                          ? formatRupiah(claim.amount)
                          : formatDuration(
                              claim.duration_hours
                            )}
                      </td>
                      <td>
                        <ClaimStatusBadge status={claim.status} />
                      </td>
                      <td>
                        {claim.reviewer_name || "Belum diperiksa"}
                      </td>
                      <td className="text-end pe-4">
                        <button
                          type="button"
                          className="btn btn-outline-primary btn-sm"
                          onClick={() =>
                            navigate(
                              `/admin/claims/${claim.claim_type}/${claim.claim_id}`
                            )
                          }
                        >
                          Detail
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="card-footer bg-white border-top p-4 d-flex flex-column flex-sm-row justify-content-between align-items-sm-center gap-3">
              <span className="text-secondary small">
                Total {state.pagination.total} data
              </span>

              <div className="d-flex align-items-center gap-2">
                <button
                  type="button"
                  className="btn btn-outline-secondary btn-sm"
                  disabled={
                    state.loading ||
                    state.pagination.page <= 1
                  }
                  onClick={() =>
                    setPage((current) => Math.max(1, current - 1))
                  }
                >
                  Sebelumnya
                </button>

                <span className="small">
                  Halaman {state.pagination.page} dari{" "}
                  {Math.max(1, state.pagination.total_pages)}
                </span>

                <button
                  type="button"
                  className="btn btn-outline-secondary btn-sm"
                  disabled={
                    state.loading ||
                    state.pagination.total_pages === 0 ||
                    state.pagination.page >=
                      state.pagination.total_pages
                  }
                  onClick={() =>
                    setPage((current) => current + 1)
                  }
                >
                  Berikutnya
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </DashboardLayout>
  );
}
