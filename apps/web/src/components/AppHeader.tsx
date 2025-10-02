import { useAuth } from "@/context/AuthContext";

export function ProjectSwitcher() {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  return (
    <label>
      Project{" "}
      <input
        aria-label="Active project id"
        defaultValue={projectId}
        onBlur={(e) => localStorage.setItem("replay_project_id", e.target.value)}
      />
    </label>
  );
}

export function UserMenu() {
  const { logout, user } = useAuth();
  return (
    <div>
      <span>{user?.email ?? "Guest"}</span>
      <button type="button" onClick={logout}>
        Logout
      </button>
    </div>
  );
}

export function AppHeader() {
  return (
    <header className="header">
      <strong>Replay</strong>
      <ProjectSwitcher />
      <UserMenu />
    </header>
  );
}
