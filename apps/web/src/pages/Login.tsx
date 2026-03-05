import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { AuthLayout } from "@/components/AuthLayout";

export function LoginForm() {
  const { login } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!email || password.length < 12) {
      setError("Email and password (12+ characters) required");
      return;
    }
    try {
      await login(email, password);
      nav("/");
    } catch {
      setError("Invalid email or password");
    }
  }

  return (
    <AuthLayout
      title="Welcome back"
      subtitle="Sign in to your Replay workspace"
      footer={
        <>
          New here? <Link to="/register">Create an account</Link>
        </>
      }
    >
      <form onSubmit={onSubmit}>
        {error && (
          <div className="alert alert--error" role="alert">
            {error}
          </div>
        )}
        <div className="form-field">
          <label htmlFor="login-email">Work email</label>
          <input
            id="login-email"
            type="email"
            autoComplete="email"
            placeholder="you@company.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <div className="form-field">
          <label htmlFor="login-password">Password</label>
          <input
            id="login-password"
            type="password"
            autoComplete="current-password"
            placeholder="••••••••••••"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit" className="btn btn--primary" style={{ width: "100%" }}>
          Sign in
        </button>
      </form>
    </AuthLayout>
  );
}

export default function Login() {
  return <LoginForm />;
}
