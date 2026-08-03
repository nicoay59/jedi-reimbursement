import { useEffect, useState } from "react";
import ClaimStatusBadge from "../components/ClaimStatusBadge";
import DashboardLayout from "../components/DashboardLayout";
import LoadingScreen from "../components/LoadingScreen";
import PageHeader from "../components/PageHeader";
import {
  getParkingClaim,
  getParkingReceipt
} from "../services/parkingClaimService";
import {
  formatDateRange,
  formatDateTime,
  formatFileSize,
  formatRupiah
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

export default function ParkingClaimDetailPage({ claimID }) {
  const [state, setState] = useState({
    loading: true,
    downloading: false,
    error: "",
    claim: null
  });

  useEffect(() => {
    let active = true;

    async function loadClaim() {
      try {
        const response = await getParkingClaim(claimID);

        if (!active) return;

        setState({
          loading: false,
          downloading: false,
          error: "",
          claim: response.data
        });
      } catch (error) {
        if (!active) return;

        setState({
          loading: false,
          downloading: false,
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

  async function downloadReceipt() {
    setState((current) => ({
      ...current,
      downloading: true,
      error: ""
    }));

    try {
      const response = await getParkingReceipt(claimID);
      const url = URL.createObjectURL(response.blob);
      const link = document.createElement("a");

      link.href = url;
      link.download =
        state.claim.receipt_original_name ||
        `bukti-parkir-${claimID}`;
      document.body.appendChild(link);
      link.click();
      link.remove();

      setTimeout(() => URL.revokeObjectURL(url), 1000);

      setState((current) => ({
        ...current,
        downloading: false
      }));
    } catch (error) {
      setState((current) => ({
        ...current,
        downloading: false,
        error: error.message
      }));
    }
  }

  if (state.loading) {
    return (
      <DashboardLayout>
        <LoadingScreen message="Memuat detail klaim..." />
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
                navigate("/employee/parking-claims")
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
        eyebrow="Detail Klaim Parkir"
        title={`Pengajuan #${claim.id}`}
        description={`Dibuat pada ${formatDateTime(
          claim.created_at
        )}`}
        actions={
          <>
            <button
              type="button"
              className="btn btn-outline-secondary"
              onClick={() =>
                navigate("/employee/parking-claims")
              }
            >
              Kembali
            </button>

            <button
              type="button"
              className="btn btn-primary"
              onClick={downloadReceipt}
              disabled={
                state.downloading ||
                !claim.receipt_available
              }
            >
              {state.downloading
                ? "Mengunduh..."
                : "Unduh bukti"}
            </button>
          </>
        }
      />

      {state.error && (
        <div className="alert alert-danger">{state.error}</div>
      )}

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
                  <p className="text-secondary mb-1">Nominal</p>
                  <div className="h4 mb-0">
                    {formatRupiah(claim.amount)}
                  </div>
                </div>
              </div>

              <dl className="row claim-detail-list mb-0">
                <dt className="col-sm-4">Periode parkir</dt>
                <dd className="col-sm-8">
                  {formatDateRange(
                    claim.parking_start_date,
                    claim.parking_end_date
                  )}
                </dd>

                <dt className="col-sm-4">Lokasi</dt>
                <dd className="col-sm-8">
                  {claim.parking_location}
                </dd>

                <dt className="col-sm-4">Deskripsi</dt>
                <dd className="col-sm-8">
                  {claim.description || "-"}
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
                Bukti pembayaran
              </p>
              <h2 className="h5 text-break">
                {claim.receipt_original_name}
              </h2>

              <dl className="row small mt-4 mb-0">
                <dt className="col-5 text-secondary">Tipe</dt>
                <dd className="col-7 text-end">
                  {claim.receipt_mime_type}
                </dd>

                <dt className="col-5 text-secondary">Ukuran</dt>
                <dd className="col-7 text-end mb-0">
                  {formatFileSize(claim.receipt_size)}
                </dd>
              </dl>
            </div>
          </aside>
        </div>
      </div>
    </DashboardLayout>
  );
}
