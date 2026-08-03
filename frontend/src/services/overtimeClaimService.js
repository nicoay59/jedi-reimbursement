import { apiRequest } from "./apiClient";

export function createOvertimeClaim(payload) {
  return apiRequest("/overtime-claims", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export function getOvertimeClaims({
  page = 1,
  limit = 10
} = {}) {
  const params = new URLSearchParams({
    page: String(page),
    limit: String(limit)
  });

  return apiRequest(`/overtime-claims?${params.toString()}`);
}

export function getOvertimeClaim(id) {
  return apiRequest(`/overtime-claims/${id}`);
}
