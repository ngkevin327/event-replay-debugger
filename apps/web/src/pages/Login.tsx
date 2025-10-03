import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

export function LoginForm() {
  const { login } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!email || password.length < 12) {
      setError("Email and password (12+ chars) required");
      return;
    }
    try {
      await login(email, password);
      nav("/");
    } catch {
      setError("Login failed");
    }
  }

  return (
    <form onSubmit={onSubmit}>
      <h1>Login</h1>
      {error && <p role="alert">{error}</p>}
      <label>
        Email
        <input value={email} onChange={(e) => setEmail(e.target.value)} />
      </label>
      <label>
        Password
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </label>
      <button type="submit">Sign in</button>
      <p>
        <Link to="/register">Register</Link>
      </p>
    </form>
  );
}

export default function Login() {
  return <LoginForm />;
}
