import useAuth from "../hooks/useAuth";
import usePathname from "../hooks/usePathname";
import { ROLES } from "../constants/roles";
import { navigate } from "../utils/navigation";
import AppIcon from "./AppIcon";

export default function DashboardNavbar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();

  async function handleLogout() {
    await logout();
    navigate("/login", { replace: true });
  }

  const dashboardPath =
    user.role === ROLES.ADMIN
      ? "/admin/dashboard"
      : "/employee/dashboard";

  const navItems =
    user.role === ROLES.ADMIN
      ? [
          { label: "Dashboard", path: dashboardPath, icon: "dashboard", exact: true },
          { label: "Pemeriksaan", path: "/admin/claims", icon: "review" },
          { label: "Laporan", path: "/admin/reports", icon: "report" }
        ]
      : [
          { label: "Dashboard", path: dashboardPath, icon: "dashboard", exact: true },
          { label: "Klaim Parkir", path: "/employee/parking-claims", icon: "parking" },
          { label: "Klaim Lembur", path: "/employee/overtime-claims", icon: "overtime" }
        ];

  function isActive(item) {
    return item.exact
      ? pathname === item.path
      : pathname.startsWith(item.path);
  }

  return (
    <nav className="navbar navbar-expand-lg sticky-top dashboard-navbar">
      <div className="container-xl navbar-inner">
        <button
          type="button"
          className="navbar-brand btn btn-link text-decoration-none d-flex align-items-center gap-3 p-0"
          onClick={() => navigate(dashboardPath)}
        >
          <span className="brand-mark brand-mark-navbar">JDI</span>
          <span className="brand-copy">
            <strong>Reimbursement Portal</strong>
            <small>PT. Jedi Global Teknologi</small>
          </span>
        </button>

        <button
          className="navbar-toggler corporate-toggler border-0 shadow-none"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#dashboardNavigation"
          aria-controls="dashboardNavigation"
          aria-expanded="false"
          aria-label="Buka navigasi"
        >
          <span className="navbar-toggler-icon" />
        </button>

        <div className="collapse navbar-collapse" id="dashboardNavigation">
          <div className="navbar-nav me-auto ms-lg-4 navigation-links">
            {navItems.map((item) => (
              <button
                key={item.path}
                type="button"
                className={`nav-link btn btn-link text-start d-flex align-items-center gap-2 ${
                  isActive(item) ? "active" : ""
                }`}
                onClick={() => navigate(item.path)}
                aria-current={isActive(item) ? "page" : undefined}
              >
                <AppIcon name={item.icon} size={17} />
                <span>{item.label}</span>
              </button>
            ))}

            {user.role === ROLES.EMPLOYEE && (
              <div className="dropdown">
                <button
                  className="nav-link dropdown-toggle btn btn-link text-start d-flex align-items-center gap-2"
                  type="button"
                  data-bs-toggle="dropdown"
                  aria-expanded="false"
                >
                  <AppIcon name="plus" size={17} />
                  <span>Pengajuan Baru</span>
                </button>

                <ul className="dropdown-menu dropdown-menu-shadow">
                  <li>
                    <button
                      type="button"
                      className="dropdown-item d-flex align-items-center gap-2"
                      onClick={() => navigate("/employee/parking-claims/new")}
                    >
                      <AppIcon name="parking" size={17} />
                      Klaim Parkir
                    </button>
                  </li>
                  <li>
                    <button
                      type="button"
                      className="dropdown-item d-flex align-items-center gap-2"
                      onClick={() => navigate("/employee/overtime-claims/new")}
                    >
                      <AppIcon name="overtime" size={17} />
                      Klaim Lembur
                    </button>
                  </li>
                </ul>
              </div>
            )}
          </div>

          <div className="d-flex align-items-center gap-3 mt-3 mt-lg-0 user-menu">
            <div className="user-avatar">{initials(user.full_name)}</div>
            <div className="d-none d-md-block user-copy">
              <div className="user-name text-truncate">{user.full_name}</div>
              <div className="dashboard-role">
                {user.role === ROLES.ADMIN ? "Administrator Sistem" : "Karyawan"}
              </div>
            </div>

            <button
              type="button"
              className="btn logout-button d-flex align-items-center gap-2"
              onClick={handleLogout}
            >
              <AppIcon name="logout" size={17} />
              <span className="d-none d-xl-inline">Keluar</span>
            </button>
          </div>
        </div>
      </div>
    </nav>
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
      .toUpperCase() || "JDI"
  );
}
