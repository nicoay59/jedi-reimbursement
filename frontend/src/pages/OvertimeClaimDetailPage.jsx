import { useEffect, useState } from "react";
import ClaimStatusBadge from "../components/ClaimStatusBadge";
import DashboardLayout from "../components/DashboardLayout";
import LoadingScreen from "../components/LoadingScreen";
import PageHeader from "../components/PageHeader";
import { getOvertimeClaim } from "../services/overtimeClaimService";
import {
  formatDate,
  formatDateTime,
  formatDuration
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

export default function OvertimeClaimDetailPage({ claimID }) {
  const [state, setState] = useState({
    loading: true,
    error: "",
    claim: null
  });

  useEffect(() => {
    let active = true;

    async function loadClaim() {
      try {
        const response = await getOvertimeClaim(claimID);

        if (!active) return;

        setState({
          loading: false,
          error: "",
          claim: response.data
        });
      } catch (error) {
        if (!active) return;

        setState({
          loading: false,
          error: error.message,
          claim: null
        });
      }
    }

    loadClaim();

    return () => {
      active = false;
    };
  }, [claimID]);

  if (state.loading) {
    return (
      <DashboardLayout>
        <LoadingScreen message="Memuat detail klaim lembur..." />
      </DashboardLayout>
    );
  }

  if (!state.claim) {
    return (
      <DashboardLayout>
        <div className="card app-card">
          <div className="card-body p-5 text-center">
            <h1 className="h4">Klaim tidak dapat dimuat</h1>
            <p className="text-secondary">{state.error}</p>
            <button
              type="button"
              className="btn btn-primary"
              onClick={() =>
                navigate("/employee/overtime-claims")
              }
            >
              Kembali ke riwayat
            </button>
          </div>
        </div>
      </DashboardLayout>
    );
  }

  const claim = state.claim;

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Detail Klaim Lembur"
        title={`Pengajuan #${claim.id}`}
        description={`Dibuat pada ${formatDateTime(
          claim.created_at
        )}`}
        actions={
          <button
            type="button"
            className="btn btn-outline-secondary"
            onClick={() =>
              navigate("/employee/overtime-claims")
            }
          >
            Kembali
          </button>
        }
      />

      <div className="row g-4">
        <div className="col-lg-8">
          <div className="card app-card">
            <div className="card-body p-4 p-lg-5">
              <div className="d-flex justify-content-between align-items-start gap-3 mb-4">
                <div>
                  <p className="text-secondary mb-1">
                    Status pengajuan
                  </p>
                  <ClaimStatusBadge status={claim.status} />
                </div>

                <div className="text-end">
                  <p className="text-secondary mb-1">Durasi</p>
                  <div className="h4 mb-0">
                    {formatDuration(claim.duration_hours)}
                  </div>
                </div>
              </div>

              <dl className="row claim-detail-list mb-0">
                <dt className="col-sm-4">Tanggal lembur</dt>
                <dd className="col-sm-8">
                  {formatDate(claim.overtime_date)}
                </dd>

                <dt className="col-sm-4">Waktu mulai</dt>
                <dd className="col-sm-8">
                  {claim.start_time}
                </dd>

                <dt className="col-sm-4">Waktu selesai</dt>
                <dd className="col-sm-8">
                  {claim.end_time}
                </dd>

                <dt className="col-sm-4">Deskripsi pekerjaan</dt>
                <dd className="col-sm-8 preserve-lines">
                  {claim.work_description}
                </dd>

                <dt className="col-sm-4">Catatan admin</dt>
                <dd className="col-sm-8">
                  {claim.admin_note || "Belum ada catatan"}
                </dd>

                <dt className="col-sm-4">Pembaruan terakhir</dt>
                <dd className="col-sm-8 mb-0">
                  {formatDateTime(claim.updated_at)}
                </dd>
              </dl>
            </div>
          </div>
        </div>

        <div className="col-lg-4">
          <aside className="card app-card">
            <div className="card-body p-4">
              <p className="text-primary fw-semibold mb-2">
                Ringkasan waktu
              </p>
              <h2 className="h5">
                {claim.start_time}–{claim.end_time}
              </h2>

              <p className="text-secondary mb-0">
                Durasi telah dihitung dan disimpan oleh backend saat
                pengajuan dibuat.
              </p>
            </div>
          </aside>
        </div>
      </div>
    </DashboardLayout>
  );
}
