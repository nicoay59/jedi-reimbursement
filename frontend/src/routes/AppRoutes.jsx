import useAuth from "../hooks/useAuth";
import usePathname from "../hooks/usePathname";
import { ROLES } from "../constants/roles";
import ProtectedRoute from "./ProtectedRoute";
import RoleRoute from "./RoleRoute";
import Redirect from "./Redirect";
import LoginPage from "../pages/LoginPage";
import AdminDashboardPage from "../pages/AdminDashboardPage";
import AdminClaimListPage from "../pages/AdminClaimListPage";
import AdminClaimDetailPage from "../pages/AdminClaimDetailPage";
import AdminReportPage from "../pages/AdminReportPage";
import EmployeeDashboardPage from "../pages/EmployeeDashboardPage";
import CreateParkingClaimPage from "../pages/CreateParkingClaimPage";
import ParkingClaimListPage from "../pages/ParkingClaimListPage";
import ParkingClaimDetailPage from "../pages/ParkingClaimDetailPage";
import CreateOvertimeClaimPage from "../pages/CreateOvertimeClaimPage";
import OvertimeClaimListPage from "../pages/OvertimeClaimListPage";
import OvertimeClaimDetailPage from "../pages/OvertimeClaimDetailPage";
import ForbiddenPage from "../pages/ForbiddenPage";
import NotFoundPage from "../pages/NotFoundPage";
import LoadingScreen from "../components/LoadingScreen";

export default function AppRoutes() {
  const pathname = usePathname();
  const { loading, authenticated, user } = useAuth();

  if (loading) {
    return <LoadingScreen message="Memulihkan sesi pengguna..." />;
  }

  if (pathname === "/login") {
    if (authenticated) {
      return <Redirect to={dashboardPath(user.role)} />;
    }
    return <LoginPage />;
  }


  if (pathname === "/admin" || pathname === "/admin/dashboard") {
    return (
      <ProtectedRoute>
        <RoleRoute roles={[ROLES.ADMIN]}>
          <AdminDashboardPage />
        </RoleRoute>
      </ProtectedRoute>
    );
  }

  if (pathname === "/admin/claims") {
    return (
      <AdminOnly>
        <AdminClaimListPage />
      </AdminOnly>
    );
  }

  if (pathname === "/admin/reports") {
    return (
      <AdminOnly>
        <AdminReportPage />
      </AdminOnly>
    );
  }

  const adminClaimMatch = pathname.match(
    /^\/admin\/claims\/(PARKING|OVERTIME)\/([1-9]\d*)$/
  );
  if (adminClaimMatch) {
    return (
      <AdminOnly>
        <AdminClaimDetailPage
          claimType={adminClaimMatch[1]}
          claimID={adminClaimMatch[2]}
        />
      </AdminOnly>
    );
  }

  if (
    pathname === "/employee" ||
    pathname === "/employee/dashboard"
  ) {
    return (
      <EmployeeOnly>
        <EmployeeDashboardPage />
      </EmployeeOnly>
    );
  }

  if (pathname === "/employee/parking-claims") {
    return (
      <EmployeeOnly>
        <ParkingClaimListPage />
      </EmployeeOnly>
    );
  }

  if (pathname === "/employee/parking-claims/new") {
    return (
      <EmployeeOnly>
        <CreateParkingClaimPage />
      </EmployeeOnly>
    );
  }

  const parkingClaimMatch = pathname.match(
    /^\/employee\/parking-claims\/([1-9]\d*)$/
  );
  if (parkingClaimMatch) {
    return (
      <EmployeeOnly>
        <ParkingClaimDetailPage
          claimID={parkingClaimMatch[1]}
        />
      </EmployeeOnly>
    );
  }

  if (pathname === "/employee/overtime-claims") {
    return (
      <EmployeeOnly>
        <OvertimeClaimListPage />
      </EmployeeOnly>
    );
  }

  if (pathname === "/employee/overtime-claims/new") {
    return (
      <EmployeeOnly>
        <CreateOvertimeClaimPage />
      </EmployeeOnly>
    );
  }

  const overtimeClaimMatch = pathname.match(
    /^\/employee\/overtime-claims\/([1-9]\d*)$/
  );
  if (overtimeClaimMatch) {
    return (
      <EmployeeOnly>
        <OvertimeClaimDetailPage
          claimID={overtimeClaimMatch[1]}
        />
      </EmployeeOnly>
    );
  }

  if (pathname === "/forbidden") {
    return <ForbiddenPage />;
  }

  if (pathname === "/") {
    return (
      <Redirect
        to={authenticated ? dashboardPath(user.role) : "/login"}
      />
    );
  }

  return <NotFoundPage />;
}

function AdminOnly({ children }) {
  return (
    <ProtectedRoute>
      <RoleRoute roles={[ROLES.ADMIN]}>
        {children}
      </RoleRoute>
    </ProtectedRoute>
  );
}

function EmployeeOnly({ children }) {
  return (
    <ProtectedRoute>
      <RoleRoute roles={[ROLES.EMPLOYEE]}>
        {children}
      </RoleRoute>
    </ProtectedRoute>
  );
}

function dashboardPath(role) {
  return role === ROLES.ADMIN
    ? "/admin/dashboard"
    : "/employee/dashboard";
}
