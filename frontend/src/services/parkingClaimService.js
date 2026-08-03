import {
  apiBlobRequest,
  apiRequest
} from "./apiClient";

export function createParkingClaim(formData) {
  return apiRequest("/parking-claims", {
    method: "POST",
    body: formData
  });
}

export function getParkingClaims({
  page = 1,
  limit = 10
} = {}) {
  const params = new URLSearchParams({
    page: String(page),
    limit: String(limit)
  });

  return apiRequest(`/parking-claims?${params.toString()}`);
}

export function getParkingClaim(id) {
  return apiRequest(`/parking-claims/${id}`);
}

export function getParkingReceipt(id) {
  return apiBlobRequest(`/parking-claims/${id}/receipt`);
}
