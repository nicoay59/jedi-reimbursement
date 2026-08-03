import DashboardLayout from "../components/DashboardLayout";
import AppIcon from "../components/AppIcon";
import useAuth from "../hooks/useAuth";
import { navigate } from "../utils/navigation";

const modules = [
  {
    icon: "parking",
    tone: "blue",
    title: "Ajukan Klaim Parkir",
    description: "Ajukan biaya parkir berdasarkan rentang tanggal dalam satu bulan.",
    action: () => navigate("/employee/parking-claims/new")
  },
  {
    icon: "history",
    tone: "slate",
    title: "Riwayat Klaim Parkir",
    description: "Pantau status, nominal, periode, dan bukti pengajuan parkir.",
    action: () => navigate("/employee/parking-claims")
  },
  {
    icon: "overtime",
    tone: "teal",
    title: "Ajukan Klaim Lembur",
    description: "Catat waktu lembur dan rincian pekerjaan secara terstruktur.",
    action: () => navigate("/employee/overtime-claims/new")
  },
  {
    icon: "history",
    tone: "purple",
    title: "Riwayat Klaim Lembur",
    description: "Lihat durasi, keputusan, dan catatan pemeriksaan klaim lembur.",
    action: () => navigate("/employee/overtime-claims")
  }
];

export default function EmployeeDashboardPage() {
  const { user } = useAuth();
  const firstName = user.full_name?.split(" ")[0] || "Karyawan";

  return (
    <DashboardLayout>
      <section className="employee-hero mb-4">
        <div className="row align-items-center g-4">
          <div className="col-lg-8">
            <span className="dashboard-eyebrow">Employee Workspace</span>
            <h1 className="employee-hero-title mt-2 mb-3">
              Selamat datang, {firstName}.
            </h1>
            <p className="employee-hero-description mb-0">
              Kelola pengajuan reimbursement dan pantau statusnya melalui satu
              portal internal yang aman dan mudah digunakan.
            </p>
          </div>

          <div className="col-lg-4">
            <div className="employee-profile-card">
              <div className="employee-avatar">{initials(user.full_name)}</div>
              <div className="min-w-0">
                <div className="fw-bold text-truncate">{user.full_name}</div>
                <div className="small text-secondary text-truncate">{user.email}</div>
                <span className="employee-role-badge mt-2">Karyawan Aktif</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="section-heading d-flex flex-column flex-sm-row justify-content-between align-items-sm-end gap-3 mb-3">
        <div>
          <p className="section-kicker mb-1">Layanan Reimbursement</p>
          <h2 className="h4 fw-bold mb-0">Pilih layanan yang dibutuhkan</h2>
        </div>
        <p className="small text-secondary mb-0">
          Seluruh aktivitas pengajuan tercatat pada sistem.
        </p>
      </section>

      <div className="row g-4">
        {modules.map((module) => (
          <div className="col-md-6" key={module.title}>
            <button
              type="button"
              className="employee-action-card w-100 text-start"
              onClick={module.action}
            >
              <span className={`employee-action-icon tone-${module.tone}`}>
                <AppIcon name={module.icon} size={24} />
              </span>
              <span className="employee-action-content">
                <strong>{module.title}</strong>
                <span>{module.description}</span>
              </span>
              <span className="employee-action-arrow" aria-hidden="true">
                <AppIcon name="arrowRight" size={22} />
              </span>
            </button>
          </div>
        ))}
      </div>

      <section className="employee-help-card mt-4">
        <div className="d-flex align-items-start gap-3">
          <span className="help-icon"><AppIcon name="info" size={22} /></span>
          <div>
            <strong>Periksa data sebelum mengirim pengajuan</strong>
            <p className="text-secondary small mb-0 mt-1">
              Pastikan periode, nominal, waktu, keterangan, dan bukti transaksi
              sudah sesuai agar proses pemeriksaan berjalan lancar.
            </p>
          </div>
        </div>
        <button
          type="button"
          className="btn btn-outline-primary"
          onClick={() => navigate("/employee/parking-claims")}
        >
          Lihat Riwayat
        </button>
      </section>
    </DashboardLayout>
  );
}

function initials(name = "") {
  return (
    name
      .split(" ")
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0])
      .join("")
      .toUpperCase() || "KR"
  );
}
