import { NavLink, Outlet } from "react-router-dom";
import { AppHeader } from "./AppHeader";

const navItems = [
  { to: "/", label: "Dashboard", end: true },
  { to: "/incidents", label: "Incidents" },
  { to: "/agents/setup", label: "Agent setup" },
  { to: "/settings", label: "Settings" },
] as const;

export function Layout() {
  return (
    <div className="layout">
      <a href="#main" className="skip-link">
        Skip to content
      </a>
      <AppHeader />
      <nav className="nav" aria-label="Main">
        {navItems.map(({ to, label, ...rest }) => (
          <NavLink key={to} to={to} {...rest}>
            {label}
          </NavLink>
        ))}
      </nav>
      <main id="main">
        <Outlet />
      </main>
    </div>
  );
}
