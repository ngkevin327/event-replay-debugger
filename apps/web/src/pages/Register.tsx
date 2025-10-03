import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

export function RegisterForm() {
  const { register } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [orgName, setOrgName] = useState("");

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    await register(email, password, orgName);
    nav("/");
  }

  return (
    <form onSubmit={onSubmit}>
      <h1>Register</h1>
      <label>
        Organization
        <input value={orgName} onChange={(e) => setOrgName(e.target.value)} />
      </label>
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
      <button type="submit">Create account</button>
      <p>
        <Link to="/login">Login</Link>
      </p>
    </form>
  );
}

export default function Register() {
  return <RegisterForm />;
}
