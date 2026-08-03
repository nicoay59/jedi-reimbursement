import {
  createContext,
  useCallback,
  useEffect,
  useMemo,
  useState
} from "react";
import {
  getCurrentUser,
  loginRequest,
  logoutRequest
} from "../services/authService";
import { TOKEN_STORAGE_KEY } from "../services/apiClient";

export const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const clearSession = useCallback(() => {
    localStorage.removeItem(TOKEN_STORAGE_KEY);
    setUser(null);
  }, []);

  const loadSession = useCallback(async () => {
    const token = localStorage.getItem(TOKEN_STORAGE_KEY);

    if (!token) {
      setLoading(false);
      return;
    }

    try {
      const response = await getCurrentUser();
      setUser(response.data);
    } catch {
      clearSession();
    } finally {
      setLoading(false);
    }
  }, [clearSession]);

  useEffect(() => {
    loadSession();
  }, [loadSession]);

  const login = useCallback(async (email, password) => {
    const response = await loginRequest({ email, password });
    const result = response.data;

    localStorage.setItem(
      TOKEN_STORAGE_KEY,
      result.access_token
    );
    setUser(result.user);

    return result.user;
  }, []);

  const logout = useCallback(async () => {
    try {
      if (localStorage.getItem(TOKEN_STORAGE_KEY)) {
        await logoutRequest();
      }
    } catch {
      // Sesi lokal tetap dihapus jika backend tidak dapat dihubungi.
    } finally {
      clearSession();
    }
  }, [clearSession]);

  const value = useMemo(
    () => ({
      user,
      loading,
      authenticated: Boolean(user),
      login,
      logout,
      refreshSession: loadSession
    }),
    [user, loading, login, logout, loadSession]
  );

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}
