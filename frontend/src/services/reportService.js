import {
  apiBlobRequest,
  apiRequest
} from "./apiClient";

export function getAdminDashboard({
  startDate,
  endDate
} = {}) {
  const params = buildPeriodParams(startDate, endDate);
  return apiRequest(`/admin/dashboard?${params.toString()}`);
}

export function getAdminReports({
  startDate,
  endDate,
  type = "ALL",
  status = "ALL",
  page = 1,
  limit = 10
} = {}) {
  const params = buildPeriodParams(startDate, endDate);
  params.set("type", type);
  params.set("status", status);
  params.set("page", String(page));
  params.set("limit", String(limit));

  return apiRequest(`/admin/reports?${params.toString()}`);
}

export function exportAdminReports({
  startDate,
  endDate,
  type = "ALL",
  status = "ALL"
} = {}) {
  const params = buildPeriodParams(startDate, endDate);
  params.set("type", type);
  params.set("status", status);

  return apiBlobRequest(
    `/admin/reports/export?${params.toString()}`
  );
}

function buildPeriodParams(startDate, endDate) {
  const params = new URLSearchParams();

  if (startDate) {
    params.set("start_date", startDate);
  }

  if (endDate) {
    params.set("end_date", endDate);
  }

  return params;
}
