import DashboardNavbar from "./DashboardNavbar";

export default function DashboardLayout({ children }) {
  return (
    <div className="app-shell min-vh-100">
      <DashboardNavbar />
      <main className="container-xl app-content">{children}</main>
      <footer className="app-footer">
        <div className="container-xl d-flex flex-column flex-sm-row justify-content-between gap-2">
          <span>Jedi Reimbursement Management System</span>
          <span>PT. Jedi Global Teknologi · Internal Use Only</span>
        </div>
      </footer>
    </div>
  );
}
