import { apiRequest } from "./apiClient";

export function getHealth() {
  return apiRequest("/health");
}
