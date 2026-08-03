import { useMemo, useState } from "react";
import DashboardLayout from "../components/DashboardLayout";
import PageHeader from "../components/PageHeader";
import AppIcon from "../components/AppIcon";
import { createParkingClaim } from "../services/parkingClaimService";
import {
  formatRupiah,
  localISODate,
  toLocalISODate
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

const MAX_CLAIM_AMOUNT = 200000;
const today = localISODate();
const minimumDate = firstDayMonthsAgo(3);

const initialForm = {
  parkingStartDate: today,
  parkingEndDate: today,
  parkingLocation: "",
  amount: "",
  description: "",
  receipt: null
};

export default function CreateParkingClaimPage() {
  const [form, setForm] = useState(initialForm);
  const [errors, setErrors] = useState({});
  const [state, setState] = useState({
    submitting: false,
    message: ""
  });

  const maximumEndDate = useMemo(
    () => endOfAllowedMonth(form.parkingStartDate),
    [form.parkingStartDate]
  );

  function updateField(event) {
    const { name, value, files } = event.target;
    setForm((current) => {
      const next = {
        ...current,
        [name]: files ? files[0] || null : value
      };

      if (name === "parkingStartDate") {
        const invalidEnd =
          !current.parkingEndDate ||
          current.parkingEndDate < value ||
          current.parkingEndDate > endOfAllowedMonth(value) ||
          current.parkingEndDate.slice(0, 7) !== value.slice(0, 7);
        if (invalidEnd) {
          next.parkingEndDate = value;
        }
      }
      return next;
    });

    setErrors((current) => ({
      ...current,
      [name]: "",
      ...(name === "parkingStartDate" ? { parking_end_date: "" } : {})
    }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const clientErrors = validate(form);
    if (Object.keys(clientErrors).length > 0) {
      setErrors(clientErrors);
      return;
    }

    const payload = new FormData();
    payload.append("parking_start_date", form.parkingStartDate);
    payload.append("parking_end_date", form.parkingEndDate);
    payload.append("parking_location", form.parkingLocation);
    payload.append("amount", form.amount);
    payload.append("description", form.description);
    payload.append("receipt", form.receipt);

    setState({ submitting: true, message: "" });
    setErrors({});

    try {
      const response = await createParkingClaim(payload);
      navigate(`/employee/parking-claims/${response.data.id}`, {
        replace: true
      });
    } catch (error) {
      setState({ submitting: false, message: error.message });
      setErrors(error.payload?.errors || {});
    }
  }

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Klaim Parkir"
        title="Ajukan klaim parkir"
        description="Rangkum transaksi parkir dalam satu bulan kalender untuk setiap pengajuan."
        actions={
          <button
            type="button"
            className="btn btn-outline-secondary"
            onClick={() => navigate("/employee/parking-claims")}
          >
            Lihat riwayat
          </button>
        }
      />

      <div className="row g-3 mb-4">
        <InfoCard
          label="Maksimal nominal"
          value={formatRupiah(MAX_CLAIM_AMOUNT)}
          helper="Untuk setiap pengajuan"
          icon="money"
        />
        <InfoCard
          label="Periode pengajuan"
          value="1 bulan"
          helper="Tanggal harus dalam bulan yang sama"
          icon="calendar"
        />
        <InfoCard
          label="Batas pengajuan"
          value="3 bulan"
          helper="Maksimal tiga bulan sebelumnya"
          icon="history"
          emphasized
        />
      </div>

      <div className="row g-4">
        <div className="col-xl-8">
          <form className="card app-card" onSubmit={handleSubmit} noValidate>
            <div className="card-body p-4 p-lg-5">
              {state.message && (
                <div className="alert alert-danger" role="alert">
                  {state.message}
                </div>
              )}

              <div className="alert alert-info border-0 mb-4">
                Satu pengajuan hanya boleh mencakup satu bulan kalender dengan
                nominal maksimal Rp200.000. Untuk bulan yang berbeda, buat
                pengajuan terpisah. Klaim dapat diajukan untuk bulan berjalan
                sampai tiga bulan sebelumnya.
              </div>

              <div className="row g-3">
                <div className="col-md-6">
                  <label htmlFor="parkingStartDate" className="form-label">
                    Tanggal mulai
                  </label>
                  <input
                    id="parkingStartDate"
                    name="parkingStartDate"
                    type="date"
                    min={minimumDate}
                    max={today}
                    className={`form-control ${
                      errors.parking_start_date ? "is-invalid" : ""
                    }`}
                    value={form.parkingStartDate}
                    onChange={updateField}
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.parking_start_date} />
                </div>

                <div className="col-md-6">
                  <label htmlFor="parkingEndDate" className="form-label">
                    Tanggal selesai
                  </label>
                  <input
                    id="parkingEndDate"
                    name="parkingEndDate"
                    type="date"
                    min={form.parkingStartDate || minimumDate}
                    max={maximumEndDate}
                    className={`form-control ${
                      errors.parking_end_date ? "is-invalid" : ""
                    }`}
                    value={form.parkingEndDate}
                    onChange={updateField}
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.parking_end_date} />
                </div>

                <div className="col-md-6">
                  <label htmlFor="amount" className="form-label">
                    Total nominal pengajuan
                  </label>
                  <div className="input-group">
                    <span className="input-group-text">Rp</span>
                    <input
                      id="amount"
                      name="amount"
                      type="number"
                      min="1"
                      max={MAX_CLAIM_AMOUNT}
                      step="1"
                      className={`form-control ${errors.amount ? "is-invalid" : ""}`}
                      value={form.amount}
                      onChange={updateField}
                      placeholder="Contoh: 150000"
                      disabled={state.submitting}
                    />
                  </div>
                  <FieldError message={errors.amount} />
                  <div className="form-text">
                    Maksimal Rp200.000 untuk setiap pengajuan.
                  </div>
                </div>

                <div className="col-md-6">
                  <label htmlFor="parkingLocation" className="form-label">
                    Lokasi / area parkir
                  </label>
                  <input
                    id="parkingLocation"
                    name="parkingLocation"
                    type="text"
                    maxLength="200"
                    className={`form-control ${
                      errors.parking_location ? "is-invalid" : ""
                    }`}
                    value={form.parkingLocation}
                    onChange={updateField}
                    placeholder="Contoh: Centennial Tower dan lokasi operasional"
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.parking_location} />
                </div>

                <div className="col-12">
                  <label htmlFor="description" className="form-label">
                    Rincian klaim
                  </label>
                  <textarea
                    id="description"
                    name="description"
                    rows="4"
                    maxLength="1000"
                    className={`form-control ${errors.description ? "is-invalid" : ""}`}
                    value={form.description}
                    onChange={updateField}
                    placeholder="Tuliskan tanggal, lokasi, atau keterangan transaksi bila diperlukan."
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.description} />
                  <div className="form-text text-end">
                    {form.description.length}/1000
                  </div>
                </div>

                <div className="col-12">
                  <label htmlFor="receipt" className="form-label">
                    Bukti pembayaran gabungan
                  </label>
                  <input
                    id="receipt"
                    name="receipt"
                    type="file"
                    accept=".jpg,.jpeg,.png,.pdf,image/jpeg,image/png,application/pdf"
                    className={`form-control ${errors.receipt ? "is-invalid" : ""}`}
                    onChange={updateField}
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.receipt} />
                  <div className="form-text">
                    Jika terdapat beberapa bukti dalam bulan yang sama, gabungkan
                    menjadi satu PDF. Format JPG, PNG, atau PDF dengan ukuran
                    maksimal 5 MB.
                  </div>

                  {form.receipt && (
                    <div className="selected-file mt-3">
                      <strong>{form.receipt.name}</strong>
                      <span className="text-secondary small">
                        {(form.receipt.size / 1024 / 1024).toFixed(2)} MB
                      </span>
                    </div>
                  )}
                </div>
              </div>
            </div>

            <div className="card-footer bg-white border-top p-4 d-flex flex-column flex-sm-row justify-content-end gap-2">
              <button
                type="button"
                className="btn btn-outline-secondary"
                onClick={() => navigate("/employee/dashboard")}
                disabled={state.submitting}
              >
                Batal
              </button>
              <button
                type="submit"
                className="btn btn-primary"
                disabled={state.submitting}
              >
                {state.submitting ? "Mengirim pengajuan..." : "Kirim pengajuan"}
              </button>
            </div>
          </form>
        </div>

        <div className="col-xl-4">
          <aside className="card app-card">
            <div className="card-body p-4">
              <p className="text-primary fw-semibold mb-2">Ketentuan klaim</p>
              <h2 className="h5">Klaim parkir per pengajuan</h2>
              <ul className="text-secondary ps-3 mb-0">
                <li className="mb-2">
                  Satu pengajuan maksimal mencakup satu bulan kalender.
                </li>
                <li className="mb-2">
                  Tanggal mulai dan selesai harus berada dalam bulan yang sama.
                </li>
                <li className="mb-2">
                  Nominal maksimal setiap pengajuan adalah Rp200.000.
                </li>
                <li className="mb-2">
                  Pengajuan bulan yang berbeda harus dibuat secara terpisah.
                </li>
                <li>
                  Karyawan dapat mengajukan bulan berjalan sampai maksimal tiga
                  bulan sebelumnya.
                </li>
              </ul>
            </div>
          </aside>
        </div>
      </div>
    </DashboardLayout>
  );
}

function InfoCard({ label, value, helper, icon, emphasized = false }) {
  return (
    <div className="col-md-4">
      <div className={`card app-card policy-summary-card h-100 ${emphasized ? "policy-summary-emphasized" : ""}`}>
        <div className="card-body p-4 d-flex align-items-start gap-3">
          <span className="policy-summary-icon"><AppIcon name={icon} size={22} /></span>
          <div>
            <div className="text-secondary small mb-2">{label}</div>
            <div className="h4 fw-bold mb-1">{value}</div>
            <div className="text-secondary small">{helper}</div>
          </div>
        </div>
      </div>
    </div>
  );
}

function FieldError({ message }) {
  if (!message) return null;
  return <div className="invalid-feedback d-block">{message}</div>;
}

function validate(form) {
  const errors = {};
  if (!form.parkingStartDate) {
    errors.parking_start_date = "Tanggal mulai wajib diisi";
  }
  if (!form.parkingEndDate) {
    errors.parking_end_date = "Tanggal selesai wajib diisi";
  } else if (form.parkingEndDate < form.parkingStartDate) {
    errors.parking_end_date = "Tanggal selesai tidak boleh sebelum tanggal mulai";
  } else if (
    form.parkingStartDate &&
    form.parkingEndDate.slice(0, 7) !== form.parkingStartDate.slice(0, 7)
  ) {
    errors.parking_end_date = "Periode harus berada dalam bulan yang sama";
  }

  if (form.parkingLocation.trim().length < 3) {
    errors.parking_location = "Lokasi parkir minimal 3 karakter";
  }

  const amount = Number(form.amount);
  if (!Number.isFinite(amount) || amount <= 0) {
    errors.amount = "Nominal harus lebih besar dari nol";
  } else if (amount > MAX_CLAIM_AMOUNT) {
    errors.amount = "Nominal maksimal setiap pengajuan adalah Rp200.000";
  }

  if (!form.receipt) {
    errors.receipt = "Bukti pembayaran wajib dipilih";
  } else if (form.receipt.size > 5 * 1024 * 1024) {
    errors.receipt = "Ukuran bukti maksimal 5 MB";
  }
  return errors;
}

function firstDayMonthsAgo(months) {
  const now = new Date();
  return toLocalISODate(new Date(now.getFullYear(), now.getMonth() - months, 1));
}

function endOfAllowedMonth(startDate) {
  if (!startDate) return today;
  const [year, month] = startDate.split("-").map(Number);
  const endOfMonth = toLocalISODate(new Date(year, month, 0));
  return endOfMonth > today ? today : endOfMonth;
}
