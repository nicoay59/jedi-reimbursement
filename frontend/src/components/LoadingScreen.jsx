export default function LoadingScreen({ message = "Memuat data..." }) {
  return (
    <main className="loading-screen min-vh-100 d-grid place-items-center">
      <div className="loading-panel text-center">
        <span className="loading-ring" aria-hidden="true" />
        <p className="fw-semibold mb-1">Mohon tunggu</p>
        <p className="text-secondary small mb-0">{message}</p>
      </div>
    </main>
  );
}
