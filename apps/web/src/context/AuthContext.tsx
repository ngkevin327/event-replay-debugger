import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { apiFetch, setOnUnauthorized } from "@/api/client";

type User = { id: string; email: string; org_id: string };

type AuthContextValue = {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, orgName: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(
    () => localStorage.getItem("replay_token"),
  );
  const [user, setUser] = useState<User | null>(() => {
    const raw = localStorage.getItem("replay_user");
    return raw ? (JSON.parse(raw) as User) : null;
  });

  const logout = useCallback(() => {
    localStorage.removeItem("replay_token");
    localStorage.removeItem("replay_user");
    setToken(null);
    setUser(null);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const res = await apiFetch<{ access_token: string; user: User }>(
      "/v1/auth/login",
      { method: "POST", body: JSON.stringify({ email, password }) },
    );
    localStorage.setItem("replay_token", res.access_token);
    localStorage.setItem("replay_user", JSON.stringify(res.user));
    setToken(res.access_token);
    setUser(res.user);
  }, []);

  const register = useCallback(
    async (email: string, password: string, orgName: string) => {
      const res = await apiFetch<{ access_token: string; user: User }>(
        "/v1/auth/register",
        {
          method: "POST",
          body: JSON.stringify({ email, password, org_name: orgName }),
        },
      );
      localStorage.setItem("replay_token", res.access_token);
      localStorage.setItem("replay_user", JSON.stringify(res.user));
      setToken(res.access_token);
      setUser(res.user);
    },
    [],
  );

  useMemo(() => {
    setOnUnauthorized(logout);
  }, [logout]);

  const value = useMemo(
    () => ({ user, token, login, register, logout }),
    [user, token, login, register, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth outside AuthProvider");
  return ctx;
}
