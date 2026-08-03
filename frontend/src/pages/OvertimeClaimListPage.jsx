import { useEffect, useState } from "react";
import ClaimStatusBadge from "../components/ClaimStatusBadge";
import DashboardLayout from "../components/DashboardLayout";
import EmptyState from "../components/EmptyState";
import LoadingScreen from "../components/LoadingScreen";
import PageHeader from "../components/PageHeader";
import { getOvertimeClaims } from "../services/overtimeClaimService";
import {
  formatDate,
  formatDuration
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

export default function OvertimeClaimListPage() {
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
        const response = await getOvertimeClaims({
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
  }, [page]);

  if (state.loading && state.items.length === 0) {
    return (
      <DashboardLayout>
        <LoadingScreen message="Memuat klaim lembur..." />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Klaim Lembur"
        title="Riwayat klaim lembur"
        description="Lihat seluruh pengajuan lembur milik akun Anda."
        actions={
          <button
            type="button"
            className="btn btn-primary"
            onClick={() =>
              navigate("/employee/overtime-claims/new")
            }
          >
            Ajukan klaim
          </button>
        }
      />

      <div className="card app-card">
        {state.error && (
          <div className="alert alert-danger m-4 mb-0">
            {state.error}
          </div>
        )}

        {state.items.length === 0 ? (
          <EmptyState
            title="Belum ada klaim lembur"
            description="Pengajuan lembur yang Anda buat akan tampil pada halaman ini."
            action={
              <button
                type="button"
                className="btn btn-primary"
                onClick={() =>
                  navigate("/employee/overtime-claims/new")
                }
              >
                Buat pengajuan pertama
              </button>
            }
          />
        ) : (
          <>
            <div className="table-responsive">
              <table className="table align-middle mb-0">
                <thead className="table-light">
                  <tr>
                    <th className="ps-4">Tanggal</th>
                    <th>Waktu</th>
                    <th>Durasi</th>
                    <th>Pekerjaan</th>
                    <th>Status</th>
                    <th className="text-end pe-4">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {state.items.map((claim) => (
                    <tr key={claim.id}>
                      <td className="ps-4">
                        {formatDate(claim.overtime_date)}
                      </td>
                      <td>
                        {claim.start_time}–{claim.end_time}
                      </td>
                      <td>
                        {formatDuration(claim.duration_hours)}
                      </td>
                      <td>
                        <div className="text-truncate overtime-description-column">
                          {claim.work_description}
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
                              `/employee/overtime-claims/${claim.id}`
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
                Total {state.pagination.total} pengajuan
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
