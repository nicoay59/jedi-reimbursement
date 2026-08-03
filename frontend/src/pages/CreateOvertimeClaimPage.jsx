import { useMemo, useState } from "react";
import DashboardLayout from "../components/DashboardLayout";
import PageHeader from "../components/PageHeader";
import { createOvertimeClaim } from "../services/overtimeClaimService";
import {
  calculateOvertimeDuration,
  localISODate
} from "../utils/formatters";
import { navigate } from "../utils/navigation";

const initialForm = {
  overtimeDate: localISODate(),
  startTime: "18:00",
  endTime: "20:00",
  workDescription: ""
};

export default function CreateOvertimeClaimPage() {
  const [form, setForm] = useState(initialForm);
  const [errors, setErrors] = useState({});
  const [state, setState] = useState({
    submitting: false,
    message: ""
  });

  const duration = useMemo(
    () =>
      calculateOvertimeDuration(
        form.startTime,
        form.endTime
      ),
    [form.startTime, form.endTime]
  );

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

  async function handleSubmit(event) {
    event.preventDefault();

    const clientErrors = validate(form, duration);
    if (Object.keys(clientErrors).length > 0) {
      setErrors(clientErrors);
      return;
    }

    setState({
      submitting: true,
      message: ""
    });
    setErrors({});

    try {
      const response = await createOvertimeClaim({
        overtime_date: form.overtimeDate,
        start_time: form.startTime,
        end_time: form.endTime,
        work_description: form.workDescription
      });

      navigate(
        `/employee/overtime-claims/${response.data.id}`,
        { replace: true }
      );
    } catch (error) {
      setState({
        submitting: false,
        message: error.message
      });
      setErrors(error.payload?.errors || {});
    }
  }

  return (
    <DashboardLayout>
      <PageHeader
        eyebrow="Klaim Lembur"
        title="Ajukan klaim lembur"
        description="Lengkapi tanggal, rentang waktu, dan pekerjaan yang dilakukan."
        actions={
          <button
            type="button"
            className="btn btn-outline-secondary"
            onClick={() =>
              navigate("/employee/overtime-claims")
            }
          >
            Lihat riwayat
          </button>
        }
      />

      <div className="row g-4">
        <div className="col-xl-8">
          <form
            className="card app-card"
            onSubmit={handleSubmit}
            noValidate
          >
            <div className="card-body p-4 p-lg-5">
              {state.message && (
                <div className="alert alert-danger" role="alert">
                  {state.message}
                </div>
              )}

              <div className="row g-3">
                <div className="col-12">
                  <label
                    htmlFor="overtimeDate"
                    className="form-label"
                  >
                    Tanggal lembur
                  </label>
                  <input
                    id="overtimeDate"
                    name="overtimeDate"
                    type="date"
                    className={`form-control ${
                      errors.overtime_date ? "is-invalid" : ""
                    }`}
                    max={localISODate()}
                    value={form.overtimeDate}
                    onChange={updateField}
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.overtime_date} />
                </div>

                <div className="col-md-6">
                  <label
                    htmlFor="startTime"
                    className="form-label"
                  >
                    Waktu mulai
                  </label>
                  <input
                    id="startTime"
                    name="startTime"
                    type="time"
                    className={`form-control ${
                      errors.start_time ? "is-invalid" : ""
                    }`}
                    value={form.startTime}
                    onChange={updateField}
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.start_time} />
                </div>

                <div className="col-md-6">
                  <label
                    htmlFor="endTime"
                    className="form-label"
                  >
                    Waktu selesai
                  </label>
                  <input
                    id="endTime"
                    name="endTime"
                    type="time"
                    className={`form-control ${
                      errors.end_time ? "is-invalid" : ""
                    }`}
                    value={form.endTime}
                    onChange={updateField}
                    disabled={state.submitting}
                  />
                  <FieldError message={errors.end_time} />
                </div>

                <div className="col-12">
                  <div
                    className={`duration-preview ${
                      duration &&
                      duration.minutes >= 30 &&
                      duration.minutes <= 16 * 60
                        ? "duration-preview-valid"
                        : "duration-preview-invalid"
                    }`}
                  >
                    <div>
                      <span className="d-block small text-secondary">
                        Estimasi durasi
                      </span>
                      <strong className="h5 mb-0">
                        {duration?.label || "-"}
                      </strong>
                    </div>

                    {duration?.crossesMidnight && (
                      <span className="badge text-bg-info">
                        Selesai hari berikutnya
                      </span>
                    )}
                  </div>
                </div>

                <div className="col-12">
                  <label
                    htmlFor="workDescription"
                    className="form-label"
                  >
                    Deskripsi pekerjaan
                  </label>
                  <textarea
                    id="workDescription"
                    name="workDescription"
                    rows="6"
                    maxLength="2000"
                    className={`form-control ${
                      errors.work_description
                        ? "is-invalid"
                        : ""
                    }`}
                    value={form.workDescription}
                    onChange={updateField}
                    placeholder="Jelaskan pekerjaan yang dilakukan selama lembur."
                    disabled={state.submitting}
                  />
                  <FieldError
                    message={errors.work_description}
                  />
                  <div className="form-text text-end">
                    {form.workDescription.length}/2000
                  </div>
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
                {state.submitting
                  ? "Mengirim pengajuan..."
                  : "Kirim pengajuan"}
              </button>
            </div>
          </form>
        </div>

        <div className="col-xl-4">
          <aside className="card app-card">
            <div className="card-body p-4">
              <p className="text-primary fw-semibold mb-2">
                Perhitungan durasi
              </p>
              <h2 className="h5">Aturan waktu lembur</h2>

              <ul className="text-secondary ps-3 mb-0">
                <li className="mb-2">
                  Durasi minimal adalah 30 menit.
                </li>
                <li className="mb-2">
                  Durasi maksimal adalah 16 jam.
                </li>
                <li className="mb-2">
                  Waktu selesai yang lebih kecil dianggap hari berikutnya.
                </li>
                <li>
                  Durasi final dihitung kembali oleh backend.
                </li>
              </ul>
            </div>
          </aside>
        </div>
      </div>
    </DashboardLayout>
  );
}

function FieldError({ message }) {
  if (!message) return null;

  return <div className="invalid-feedback d-block">{message}</div>;
}

function validate(form, duration) {
  const errors = {};

  if (!form.overtimeDate) {
    errors.overtime_date = "Tanggal lembur wajib diisi";
  }

  if (!form.startTime) {
    errors.start_time = "Waktu mulai wajib diisi";
  }

  if (!form.endTime) {
    errors.end_time = "Waktu selesai wajib diisi";
  }

  if (duration && duration.minutes < 30) {
    errors.end_time = "Durasi lembur minimal 30 menit";
  }

  if (duration && duration.minutes > 16 * 60) {
    errors.end_time = "Durasi lembur maksimal 16 jam";
  }

  const description = form.workDescription.trim();
  if (description.length < 10) {
    errors.work_description =
      "Deskripsi pekerjaan minimal 10 karakter";
  }

  return errors;
}
