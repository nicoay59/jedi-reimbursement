const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ||
  "http://localhost:8080/api/v1";

const TOKEN_STORAGE_KEY = "jedi_reimbursement_access_token";

export async function apiRequest(path, options = {}) {
  return request(path, options, "json");
}

export async function apiBlobRequest(path, options = {}) {
  return request(path, options, "blob");
}

async function request(path, options, responseType) {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 15000);
  const token = localStorage.getItem(TOKEN_STORAGE_KEY);
  const bodyIsFormData =
    typeof FormData !== "undefined" &&
    options.body instanceof FormData;

  const headers = {
    Accept:
      responseType === "blob"
        ? "*/*"
        : "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...(!bodyIsFormData && options.body
      ? { "Content-Type": "application/json" }
      : {}),
    ...options.headers
  };

  try {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      ...options,
      signal: controller.signal,
      headers
    });

    if (!response.ok) {
      const payload = await response.json().catch(() => null);
      const error = new Error(
        payload?.message ||
          `Permintaan gagal dengan status ${response.status}`
      );
      error.status = response.status;
      error.payload = payload;
      throw error;
    }

    if (responseType === "blob") {
      return {
        blob: await response.blob(),
        contentType:
          response.headers.get("Content-Type") ||
          "application/octet-stream",
        disposition:
          response.headers.get("Content-Disposition") || ""
      };
    }

    return response.json();
  } catch (error) {
    if (error?.name === "AbortError") {
      throw new Error("Waktu koneksi ke backend habis");
    }

    if (error instanceof TypeError) {
      throw new Error(
        "Backend belum berjalan atau alamat API tidak sesuai"
      );
    }

    throw error;
  } finally {
    clearTimeout(timeout);
  }
}

export { TOKEN_STORAGE_KEY };
