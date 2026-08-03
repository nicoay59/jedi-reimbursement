import { useCallback, useEffect, useState } from "react";
import ClaimStatusBadge from "../components/ClaimStatusBadge";
import ClaimTypeBadge from "../components/ClaimTypeBadge";
import DashboardLayout from "../components/DashboardLayout";
import LoadingScreen from "../components/LoadingScreen";
import PageHeader from "../components/PageHeader";
import {
  getAdminClaim,
  getAdminClaimHistory,
  getAdminClaimReceipt,
  reviewAdminClaim
} from "../services/adminClaimService";
import {
  formatDate,
  formatDateRange,
  formatDateTime,
  formatDuration,
  formatFileSize,
  formatRupiah
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

export default function AdminClaimDetailPage({
  claimType,
  claimID
}) {
  const [state, setState] = useState({
    loading: true,
    reviewing: false,
    downloading: false,
    error: "",
    success: "",
    claim: null,
    history: []
  });
  const [form, setForm] = useState({
    status: "APPROVED",
    note: ""
  });
  const [errors, setErrors] = useState({});

  const loadData = useCallback(async () => {
    setState((current) => ({
      ...current,
      loading: true,
      error: ""
    }));

    try {
      const claimResponse = await getAdminClaim(
        claimType,
        claimID
      );

      setState((current) => ({
        ...current,
        loading: false,
        claim: claimResponse.data,
        error: ""
      }));

      try {
        const historyResponse = await getAdminClaimHistory(
          claimType,
          claimID
        );

        setState((current) => ({
          ...current,
          history: historyResponse.data.items
        }));
      } catch (historyError) {
        setState((current) => ({
          ...current,
          error: `Detail berhasil dimuat, tetapi riwayat gagal: ${
            historyError.message
          }`
        }));
      }
    } catch (error) {
      setState((current) => ({
        ...current,
        loading: false,
        claim: null,
        error: error.message
      }));
    }
  }, [claimType, claimID]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  function updateField(event) {
    const { name, value } = event.target;
    setForm((current) => ({
      ...current,
      [name]: value
    }));
    setErrors((current) => ({
      ...current,
      [name]: ""
    }));
  }

  async function handleReview(event) {
    event.preventDefault();

    const validationErrors = {};
    if (form.status === "REJECTED" && form.note.trim().length < 5) {
      validationErrors.note =
        "Catatan penolakan minimal 5 karakter";
    }
    if (form.note.length > 1000) {
      validationErrors.note =
        "Catatan maksimal 1000 karakter";
    }

    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    setState((current) => ({
      ...current,
      reviewing: true,
      error: "",
      success: ""
    }));
    setErrors({});

    try {
      const response = await reviewAdminClaim(
        claimType,
        claimID,
        {
          status: form.status,
          note: form.note
        }
      );

      setState((current) => ({
        ...current,
        reviewing: false,
        claim: response.data,
        success:
          form.status === "APPROVED"
            ? "Klaim berhasil disetujui."
            : "Klaim berhasil ditolak."
      }));

      try {
        const historyResponse = await getAdminClaimHistory(
          claimType,
          claimID
        );

        setState((current) => ({
          ...current,
          history: historyResponse.data.items
        }));
      } catch (historyError) {
        setState((current) => ({
          ...current,
          error: `Keputusan tersimpan, tetapi riwayat gagal dimuat: ${
            historyError.message
          }`
        }));
      }
    } catch (error) {
      setState((current) => ({
        ...current,
        reviewing: false,
        error: error.message
      }));
      setErrors(error.payload?.errors || {});
    }
  }

  async function downloadReceipt() {
    setState((current) => ({
      ...current,
      downloading: true,
      error: ""
    }));

    try {
      const response = await getAdminClaimReceipt(
        claimType,
        claimID
      );
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
        <LoadingScreen message="Memuat detail pemeriksaan..." />
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
              onClick={() => navigate("/admin/claims")}
            >
              Kembali ke pemeriksaan
            </button>
          </div>
        </div>
      </DashboardLayout>
    );
  }

  const claim = state.claim;
  const isPending = claim.status === "PENDING";
  const isParking = claim.claim_type === "PARKING";

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Pemeriksaan Klaim"
        title={`${isParking ? "Klaim Parkir" : "Klaim Lembur"} #${
          claim.claim_id
        }`}
        description={`Diajukan oleh ${claim.employee_name}`}
        actions={
          <>
            <button
              type="button"
              className="btn btn-outline-secondary"
              onClick={() => navigate("/admin/claims")}
            >
              Kembali
            </button>

            {isParking && claim.receipt_available && (
              <button
                type="button"
                className="btn btn-outline-primary"
                onClick={downloadReceipt}
                disabled={state.downloading}
              >
                {state.downloading
                  ? "Mengunduh..."
                  : "Unduh bukti"}
              </button>
            )}
          </>
        }
      />

      {state.error && (
        <div className="alert alert-danger">{state.error}</div>
      )}
      {state.success && (
        <div className="alert alert-success">{state.success}</div>
      )}

      <div className="row g-4">
        <div className="col-xl-8">
          <div className="card app-card mb-4">
            <div className="card-body p-4 p-lg-5">
              <div className="d-flex flex-wrap justify-content-between align-items-start gap-3 mb-4">
                <div className="d-flex gap-2">
                  <ClaimTypeBadge type={claim.claim_type} />
                  <ClaimStatusBadge status={claim.status} />
                </div>

                <div className="text-end">
                  <p className="text-secondary mb-1">
                    {isParking ? "Nominal" : "Durasi"}
                  </p>
                  <div className="h4 mb-0">
                    {isParking
                      ? formatRupiah(claim.amount)
                      : formatDuration(claim.duration_hours)}
                  </div>
                </div>
              </div>

              <dl className="row claim-detail-list mb-0">
                <dt className="col-sm-4">Karyawan</dt>
                <dd className="col-sm-8">
                  {claim.employee_name} ({claim.employee_number})
                </dd>

                <dt className="col-sm-4">
                  {isParking ? "Tanggal parkir" : "Tanggal lembur"}
                </dt>
                <dd className="col-sm-8">
                  {isParking
                    ? formatDateRange(claim.claim_date, claim.claim_end_date)
                    : formatDate(claim.claim_date)}
                </dd>

                {isParking ? (
                  <>
                    <dt className="col-sm-4">Lokasi</dt>
                    <dd className="col-sm-8">{claim.title}</dd>
                  </>
                ) : (
                  <>
                    <dt className="col-sm-4">Waktu</dt>
                    <dd className="col-sm-8">
                      {claim.start_time}–{claim.end_time}
                    </dd>
                  </>
                )}

                <dt className="col-sm-4">
                  {isParking ? "Deskripsi" : "Pekerjaan"}
                </dt>
                <dd className="col-sm-8 preserve-lines">
                  {claim.description || "-"}
                </dd>

                <dt className="col-sm-4">Diajukan</dt>
                <dd className="col-sm-8">
                  {formatDateTime(claim.created_at)}
                </dd>

                <dt className="col-sm-4">Pemeriksa</dt>
                <dd className="col-sm-8">
                  {claim.reviewer_name || "Belum diperiksa"}
                </dd>

                <dt className="col-sm-4">Catatan admin</dt>
                <dd className="col-sm-8 mb-0 preserve-lines">
                  {claim.admin_note || "Belum ada catatan"}
                </dd>
              </dl>
            </div>
          </div>

          <div className="card app-card">
            <div className="card-body p-4">
              <p className="text-primary fw-semibold mb-2">
                Riwayat status
              </p>
              <h2 className="h5 mb-4">Perjalanan pengajuan</h2>

              {state.history.length === 0 ? (
                <p className="text-secondary mb-0">
                  Riwayat belum tersedia.
                </p>
              ) : (
                <div className="claim-timeline">
                  {state.history.map((history) => (
                    <div
                      className="claim-timeline-item"
                      key={history.id}
                    >
                      <div className="claim-timeline-marker" />
                      <div className="claim-timeline-content">
                        <div className="d-flex flex-wrap justify-content-between gap-2">
                          <div>
                            <ClaimStatusBadge
                              status={history.new_status}
                            />
                            <span className="ms-2 fw-semibold">
                              {history.updated_by_name}
                            </span>
                          </div>
                          <span className="small text-secondary">
                            {formatDateTime(history.created_at)}
                          </span>
                        </div>
                        <p className="text-secondary mt-2 mb-0">
                          {history.note || "Tanpa catatan"}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="col-xl-4">
          {isPending ? (
            <form
              className="card app-card sticky-xl-top review-form-card"
              onSubmit={handleReview}
            >
              <div className="card-body p-4">
                <p className="text-primary fw-semibold mb-2">
                  Keputusan
                </p>
                <h2 className="h5">Proses pengajuan</h2>

                <div className="mb-3 mt-4">
                  <label htmlFor="status" className="form-label">
                    Hasil pemeriksaan
                  </label>
                  <select
                    id="status"
                    name="status"
                    className={`form-select ${
                      errors.status ? "is-invalid" : ""
                    }`}
                    value={form.status}
                    onChange={updateField}
                    disabled={state.reviewing}
                  >
                    <option value="APPROVED">
                      Setujui pengajuan
                    </option>
                    <option value="REJECTED">
                      Tolak pengajuan
                    </option>
                  </select>
                  <FieldError message={errors.status} />
                </div>

                <div className="mb-3">
                  <label htmlFor="note" className="form-label">
                    Catatan administrator
                  </label>
                  <textarea
                    id="note"
                    name="note"
                    rows="5"
                    maxLength="1000"
                    className={`form-control ${
                      errors.note ? "is-invalid" : ""
                    }`}
                    value={form.note}
                    onChange={updateField}
                    placeholder={
                      form.status === "REJECTED"
                        ? "Jelaskan alasan penolakan."
                        : "Catatan opsional untuk karyawan."
                    }
                    disabled={state.reviewing}
                  />
                  <FieldError message={errors.note} />
                  <div className="form-text text-end">
                    {form.note.length}/1000
                  </div>
                </div>
              </div>

              <div className="card-footer bg-white p-4">
                <button
                  type="submit"
                  className={`btn w-100 ${
                    form.status === "APPROVED"
                      ? "btn-success"
                      : "btn-danger"
                  }`}
                  disabled={state.reviewing}
                >
                  {state.reviewing
                    ? "Memproses..."
                    : form.status === "APPROVED"
                      ? "Setujui klaim"
                      : "Tolak klaim"}
                </button>
              </div>
            </form>
          ) : (
            <aside className="card app-card">
              <div className="card-body p-4">
                <p className="text-primary fw-semibold mb-2">
                  Hasil pemeriksaan
                </p>
                <h2 className="h5 mb-3">
                  Pengajuan telah diputuskan
                </h2>

                <ClaimStatusBadge status={claim.status} />

                <dl className="row small mt-4 mb-0">
                  <dt className="col-5 text-secondary">
                    Pemeriksa
                  </dt>
                  <dd className="col-7 text-end">
                    {claim.reviewer_name || "-"}
                  </dd>

                  <dt className="col-5 text-secondary">Waktu</dt>
                  <dd className="col-7 text-end mb-0">
                    {formatDateTime(claim.reviewed_at)}
                  </dd>
                </dl>
              </div>
            </aside>
          )}

          {isParking && claim.receipt_available && (
            <aside className="card app-card mt-4">
              <div className="card-body p-4">
                <p className="text-primary fw-semibold mb-2">
                  Bukti pembayaran
                </p>
                <h2 className="h6 text-break">
                  {claim.receipt_original_name}
                </h2>
                <p className="small text-secondary mb-0">
                  {claim.receipt_mime_type} ·{" "}
                  {formatFileSize(claim.receipt_size)}
                </p>
              </div>
            </aside>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}

function FieldError({ message }) {
  if (!message) return null;
  return <div className="invalid-feedback d-block">{message}</div>;
}
