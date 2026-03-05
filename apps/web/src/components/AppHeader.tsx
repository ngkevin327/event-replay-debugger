import { useAuth } from "@/context/AuthContext";
import { BrandLogo } from "./BrandLogo";
import { ThemeToggle } from "./ThemeToggle";

export function ProjectSwitcher() {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  return (
    <div className="header__project">
      <label htmlFor="project-id">Project</label>
      <input
        id="project-id"
        aria-label="Active project id"
        placeholder="Project UUID"
        defaultValue={projectId}
        onBlur={(e) => localStorage.setItem("replay_project_id", e.target.value)}
      />
    </div>
  );
}

export function UserMenu() {
  const { logout, user } = useAuth();
  return (
    <div className="header__user">
      <span title={user?.email}>{user?.email ?? "Guest"}</span>
      <button type="button" className="btn btn--ghost" onClick={logout}>
        Logout
      </button>
    </div>
  );
}

export function AppHeader() {
  return (
    <header className="header">
      <BrandLogo size="sm" />
      <div className="header__spacer" />
      <ProjectSwitcher />
      <ThemeToggle />
      <UserMenu />
    </header>
  );
}
