import { apiRequest } from "./apiClient";

export function checkAdminAccess() {
  return apiRequest("/admin/ping");
}

export function checkEmployeeAccess() {
  return apiRequest("/employee/ping");
}
