import { useState } from "react";
import useAuth from "../hooks/useAuth";
import { navigate } from "../utils/navigation";
import { ROLES } from "../constants/roles";
import AppIcon from "../components/AppIcon";

export default function LoginPage() {
  const { login } = useAuth();
  const [form, setForm] = useState({ email: "", password: "" });
  const [showPassword, setShowPassword] = useState(false);
  const [state, setState] = useState({ submitting: false, error: "" });

  function updateField(event) {
    const { name, value } = event.target;
    setForm((current) => ({ ...current, [name]: value }));
    if (state.error) setState((current) => ({ ...current, error: "" }));
  }

  async function handleSubmit(event) {
    event.preventDefault();

    if (!form.email.trim() || !form.password) {
      setState({ submitting: false, error: "Email dan password wajib diisi." });
      return;
    }

    setState({ submitting: true, error: "" });

    try {
      const user = await login(form.email.trim(), form.password);
      navigate(
        user.role === ROLES.ADMIN ? "/admin/dashboard" : "/employee/dashboard",
        { replace: true }
      );
    } catch (error) {
      setState({ submitting: false, error: error.message });
    }
  }

  return (
    <main className="login-page min-vh-100">
      <div className="container login-container">
        <div className="login-shell row g-0">
          <section className="col-lg-6 login-intro p-4 p-md-5">
            <div className="login-intro-content">
              <div className="login-brand">
                <span className="brand-mark brand-mark-lg">JDI</span>
                <div>
                  <strong className="d-block">Reimbursement Portal</strong>
                  <span className="small opacity-75">PT. Jedi Global Teknologi</span>
                </div>
              </div>

              <div className="login-copy">
                <span className="login-kicker">Enterprise Internal System</span>
                <h1 className="display-5 fw-bold mt-3 mb-3">
                  Kelola reimbursement secara terpusat dan akuntabel.
                </h1>
                <p className="lead mb-4">
                  Portal internal untuk pengajuan klaim parkir dan lembur,
                  pemeriksaan administrator, serta pelaporan yang lebih tertib.
                </p>

                <div className="login-benefits">
                  <Benefit icon="claims" text="Pengajuan dan dokumen tersimpan terpusat" />
                  <Benefit icon="history" text="Status dan riwayat dapat dipantau" />
                  <Benefit icon="shield" text="Akses aman berdasarkan peran pengguna" />
                </div>
              </div>

              <div className="login-footer-note d-flex align-items-center gap-2">
                <span className="security-dot" />
                <span>Portal internal perusahaan · Authorized personnel only</span>
              </div>
            </div>
          </section>

          <section className="col-lg-6 login-form-panel p-4 p-md-5">
            <div className="login-form-wrap">
              <div className="login-form-heading mb-4">
                <span className="login-kicker text-primary">Secure Access</span>
                <h2 className="h2 fw-bold mt-2 mb-2">Masuk ke portal</h2>
                <p className="text-secondary mb-0">
                  Gunakan akun perusahaan yang sudah terdaftar pada sistem.
                </p>
              </div>

              {state.error && (
                <div className="alert alert-danger login-alert d-flex gap-2" role="alert">
                  <AppIcon name="info" size={19} className="flex-shrink-0 mt-1" />
                  <span>{state.error}</span>
                </div>
              )}

              <form onSubmit={handleSubmit} noValidate>
                <div className="mb-3">
                  <label htmlFor="email" className="form-label fw-semibold">
                    Email perusahaan
                  </label>
                  <div className="input-group input-group-lg app-input-group">
                    <span className="input-group-text" aria-hidden="true">
                      <AppIcon name="mail" size={19} />
                    </span>
                    <input
                      id="email"
                      name="email"
                      type="email"
                      className="form-control"
                      autoComplete="email"
                      value={form.email}
                      onChange={updateField}
                      placeholder="nama@perusahaan.com"
                      disabled={state.submitting}
                      autoFocus
                    />
                  </div>
                </div>

                <div className="mb-4">
                  <label htmlFor="password" className="form-label fw-semibold">
                    Password
                  </label>
                  <div className="input-group input-group-lg app-input-group">
                    <span className="input-group-text" aria-hidden="true">
                      <AppIcon name="lock" size={19} />
                    </span>
                    <input
                      id="password"
                      name="password"
                      type={showPassword ? "text" : "password"}
                      className="form-control"
                      autoComplete="current-password"
                      value={form.password}
                      onChange={updateField}
                      placeholder="Masukkan password"
                      disabled={state.submitting}
                    />
                    <button
                      type="button"
                      className="btn password-toggle"
                      onClick={() => setShowPassword((current) => !current)}
                      disabled={state.submitting}
                      aria-label={showPassword ? "Sembunyikan password" : "Tampilkan password"}
                    >
                      <AppIcon name={showPassword ? "eyeOff" : "eye"} size={18} />
                      <span>{showPassword ? "Sembunyikan" : "Tampilkan"}</span>
                    </button>
                  </div>
                </div>

                <button
                  type="submit"
                  className="btn btn-primary btn-lg w-100 login-submit"
                  disabled={state.submitting}
                >
                  {state.submitting ? (
                    <>
                      <span className="spinner-border spinner-border-sm me-2" aria-hidden="true" />
                      Memverifikasi akun...
                    </>
                  ) : (
                    <span className="d-flex justify-content-center align-items-center gap-2">
                      Masuk ke Portal
                      <AppIcon name="arrowRight" size={18} />
                    </span>
                  )}
                </button>
              </form>

              <div className="login-support mt-4">
                <AppIcon name="shield" size={18} />
                <span>
                  Hubungi administrator sistem apabila mengalami kendala akses.
                </span>
              </div>
            </div>
          </section>
        </div>
      </div>
    </main>
  );
}

function Benefit({ icon, text }) {
  return (
    <div className="login-benefit">
      <span className="login-benefit-check" aria-hidden="true">
        <AppIcon name={icon} size={16} />
      </span>
      <span>{text}</span>
    </div>
  );
}
