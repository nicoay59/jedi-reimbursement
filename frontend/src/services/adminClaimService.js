import {
  apiBlobRequest,
  apiRequest
} from "./apiClient";

export function getAdminClaims({
  type = "ALL",
  status = "ALL",
  page = 1,
  limit = 10
} = {}) {
  const params = new URLSearchParams({
    type,
    status,
    page: String(page),
    limit: String(limit)
  });

  return apiRequest(`/admin/claims?${params.toString()}`);
}

export function getAdminClaim(type, id) {
  return apiRequest(`/admin/claims/${type}/${id}`);
}

export function reviewAdminClaim(type, id, payload) {
  return apiRequest(`/admin/claims/${type}/${id}/status`, {
    method: "PATCH",
    body: JSON.stringify(payload)
  });
}

export function getAdminClaimHistory(type, id) {
  return apiRequest(`/admin/claims/${type}/${id}/history`);
}

export function getAdminClaimReceipt(type, id) {
  return apiBlobRequest(`/admin/claims/${type}/${id}/receipt`);
}
