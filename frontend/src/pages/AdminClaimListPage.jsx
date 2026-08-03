import { useEffect, useState } from "react";
import ClaimStatusBadge from "../components/ClaimStatusBadge";
import ClaimTypeBadge from "../components/ClaimTypeBadge";
import DashboardLayout from "../components/DashboardLayout";
import EmptyState from "../components/EmptyState";
import LoadingScreen from "../components/LoadingScreen";
import PageHeader from "../components/PageHeader";
import { getAdminClaims } from "../services/adminClaimService";
import {
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

export default function AdminClaimListPage() {
  const [filters, setFilters] = useState({
    type: "ALL",
    status: "PENDING"
  });
  const [page, setPage] = useState(1);
  const [state, setState] = useState({
    loading: true,
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

    async function loadClaims() {
      setState((current) => ({
        ...current,
        loading: true,
        error: ""
      }));

      try {
        const response = await getAdminClaims({
          ...filters,
          page,
          limit: 10
        });

        if (!active) return;

        setState({
          loading: false,
          error: "",
          items: response.data.items,
          pagination: response.data.pagination
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

    loadClaims();

    return () => {
      active = false;
    };
  }, [filters, page]);

  function updateFilter(event) {
    const { name, value } = event.target;
    setFilters((current) => ({
      ...current,
      [name]: value
    }));
    setPage(1);
  }

  if (state.loading && state.items.length === 0) {
    return (
      <DashboardLayout>
        <LoadingScreen message="Memuat daftar pemeriksaan..." />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Administrator"
        title="Pemeriksaan klaim"
        description="Periksa pengajuan parkir dan lembur dari seluruh karyawan."
      />

      <div className="card app-card mb-4">
        <div className="card-body p-4">
          <div className="row g-3">
            <div className="col-md-6">
              <label htmlFor="type" className="form-label">
                Jenis klaim
              </label>
              <select
                id="type"
                name="type"
                className="form-select"
                value={filters.type}
                onChange={updateFilter}
              >
                {typeOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            <div className="col-md-6">
              <label htmlFor="status" className="form-label">
                Status
              </label>
              <select
                id="status"
                name="status"
                className="form-select"
                value={filters.status}
                onChange={updateFilter}
              >
                {statusOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
      </div>

      <div className="card app-card">
        {state.error && (
          <div className="alert alert-danger m-4 mb-0">
            {state.error}
          </div>
        )}

        {state.items.length === 0 ? (
          <EmptyState
            title="Tidak ada klaim"
            description="Tidak ditemukan klaim yang sesuai dengan filter."
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
                    <th>Ringkasan</th>
                    <th>Status</th>
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
                        <div className="fw-semibold">
                          {claim.claim_type === "PARKING"
                            ? formatRupiah(claim.amount)
                            : formatDuration(
                                claim.duration_hours
                              )}
                        </div>
                        <div className="small text-secondary text-truncate admin-claim-summary">
                          {claim.claim_type === "PARKING"
                            ? claim.title
                            : `${claim.start_time}–${claim.end_time}`}
                        </div>
                      </td>
                      <td>
                        <ClaimStatusBadge status={claim.status} />
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
                          Periksa
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="card-footer bg-white border-top p-4 d-flex flex-column flex-sm-row justify-content-between align-items-sm-center gap-3">
              <span className="text-secondary small">
                Total {state.pagination.total} klaim
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
