import { navigate } from "../utils/navigation";

export default function NotFoundPage() {
  return (
    <main className="min-vh-100 d-grid place-items-center p-4">
      <div className="text-center">
        <div className="error-code">404</div>
        <h1 className="h2">Halaman tidak ditemukan</h1>
        <p className="text-secondary">
          Alamat yang dibuka tidak tersedia.
        </p>
        <button
          type="button"
          className="btn btn-primary"
          onClick={() => navigate("/", { replace: true })}
        >
          Kembali ke beranda
        </button>
      </div>
    </main>
  );
}
