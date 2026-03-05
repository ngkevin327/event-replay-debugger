import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { AuthLayout } from "@/components/AuthLayout";

export function RegisterForm() {
  const { register } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [orgName, setOrgName] = useState("");
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!orgName.trim() || !email || password.length < 12) {
      setError("Organization, email, and password (12+ characters) required");
      return;
    }
    try {
      await register(email, password, orgName);
      nav("/");
    } catch {
      setError("Could not create account. Try a different email.");
    }
  }

  return (
    <AuthLayout
      title="Start debugging smarter"
      subtitle="Create your organization workspace"
      footer={
        <>
          Already have an account? <Link to="/login">Sign in</Link>
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
          <label htmlFor="reg-org">Organization</label>
          <input
            id="reg-org"
            placeholder="Acme Payments"
            value={orgName}
            onChange={(e) => setOrgName(e.target.value)}
          />
        </div>
        <div className="form-field">
          <label htmlFor="reg-email">Work email</label>
          <input
            id="reg-email"
            type="email"
            autoComplete="email"
            placeholder="you@company.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <div className="form-field">
          <label htmlFor="reg-password">Password</label>
          <input
            id="reg-password"
            type="password"
            autoComplete="new-password"
            placeholder="Minimum 12 characters"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit" className="btn btn--primary" style={{ width: "100%" }}>
          Create account
        </button>
      </form>
    </AuthLayout>
  );
}

export default function Register() {
  return <RegisterForm />;
}
