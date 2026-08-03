import useAuth from "../hooks/useAuth";
import { navigate } from "../utils/navigation";
import { ROLES } from "../constants/roles";

export default function ForbiddenPage() {
  const { user, authenticated } = useAuth();

  function goBack() {
    if (!authenticated) {
      navigate("/login", { replace: true });
      return;
    }

    navigate(
      user.role === ROLES.ADMIN
        ? "/admin/dashboard"
        : "/employee/dashboard",
      { replace: true }
    );
  }

  return (
    <main className="min-vh-100 d-grid place-items-center p-4">
      <div className="text-center">
        <div className="error-code">403</div>
        <h1 className="h2">Akses ditolak</h1>
        <p className="text-secondary">
          Anda tidak memiliki role yang diperlukan untuk halaman ini.
        </p>
        <button type="button" className="btn btn-primary" onClick={goBack}>
          Kembali
        </button>
      </div>
    </main>
  );
}
